//! Offline tests against a fake server.

use std::collections::{BTreeMap, BTreeSet};
use std::future::Future;
use std::pin::Pin;
use std::time::Duration;

use open_alldebrid::models::PinGetData;
use open_alldebrid::{AllDebrid, Error, TorrentFile, error_codes};
use wiremock::matchers::any;
use wiremock::{Mock, MockServer, Request, ResponseTemplate};

const SPEC: &str = concat!(env!("CARGO_MANIFEST_DIR"), "/../../openapi.yaml");
const USER: &str = r#"{"user":{"username":"u","email":"e","isPremium":true,"isSubscribed":false,"isTrial":false,"premiumUntil":0,"lang":"en","preferedDomain":"com","fidelityPoints":0,"limitedHostersQuotas":[],"notifications":[]}}"#;

struct Fake {
    server: MockServer,
    client: AllDebrid,
}

impl Fake {
    async fn new() -> Self {
        let server = MockServer::start().await;
        let client = AllDebrid::new("test-key").base_url(server.uri()).user_agent("unit");
        Self { server, client }
    }

    /// Queues one answer; answers are served in order.
    async fn answer(&self, status: u16, body: &str) {
        Mock::given(any())
            .respond_with(ResponseTemplate::new(status).set_body_string(body.to_owned()))
            .up_to_n_times(1)
            .mount(&self.server)
            .await;
    }

    async fn ok(&self, data: &str) {
        self.answer(200, &format!(r#"{{"status":"success","data":{data}}}"#)).await;
    }

    async fn requests(&self) -> Vec<Request> {
        self.server.received_requests().await.unwrap_or_default()
    }

    async fn last(&self) -> Request {
        self.requests().await.pop().expect("no request")
    }
}

fn form(req: &Request) -> BTreeMap<String, Vec<String>> {
    let mut out: BTreeMap<String, Vec<String>> = BTreeMap::new();
    for (k, v) in url::form_urlencoded::parse(&req.body) {
        out.entry(k.into_owned()).or_default().push(v.into_owned());
    }
    out
}

fn expect_form(pairs: &[(&str, &[&str])]) -> BTreeMap<String, Vec<String>> {
    pairs.iter().map(|(k, v)| ((*k).to_owned(), v.iter().map(|s| (*s).to_owned()).collect())).collect()
}

fn spec_operations() -> BTreeSet<String> {
    let text = std::fs::read_to_string(SPEC).expect("read spec");
    let mut ops = BTreeSet::new();
    let mut path = None;
    for line in text.lines() {
        if line == "components:" {
            break;
        }
        if let Some(p) = line.strip_prefix("  /").and_then(|l| l.strip_suffix(':')) {
            path = Some(format!("/{p}"));
        } else if let (Some(p), Some(m)) = (&path, line.strip_prefix("    ").and_then(|l| l.strip_suffix(':'))) {
            if m == "get" || m == "post" {
                ops.insert(format!("{} {p}", m.to_uppercase()));
            }
        }
    }
    ops
}

#[tokio::test]
async fn list_parameters_are_php_style() {
    let f = Fake::new().await;
    f.ok(r#"{"magnets":[{"magnet":"a","id":1},{"magnet":"b","error":{"code":"MAGNET_INVALID_URI","message":"m"}}]}"#)
        .await;
    let result = f.client.upload_magnets(&["a", "b"]).await.unwrap();
    let req = f.last().await;
    assert_eq!(req.method.as_str(), "POST");
    assert_eq!(req.url.path(), "/v4/magnet/upload");
    assert_eq!(form(&req), expect_form(&[("magnets[]", &["a", "b"])]));
    assert_eq!(req.headers.get("authorization").unwrap(), "Bearer test-key");
    assert_eq!(req.headers.get("user-agent").unwrap(), "unit");
    assert_eq!(result[0].id, Some(1));
    assert_eq!(result[1].error.as_ref().unwrap().code, "MAGNET_INVALID_URI");
}

#[tokio::test]
async fn no_key_no_authorization() {
    let f = Fake::new().await;
    f.ok(r#"{"ping":"pong"}"#).await;
    f.client.with_api_key(None).ping().await.unwrap();
    assert!(f.last().await.headers.get("authorization").is_none());
}

#[tokio::test]
async fn get_parameters_in_query() {
    let f = Fake::new().await;
    f.ok(r#"{"hosts":{}}"#).await;
    f.client.hosts(true).await.unwrap();
    let req = f.last().await;
    assert_eq!(req.method.as_str(), "GET");
    assert_eq!(req.url.query(), Some("hostsOnly=true"));
}

#[tokio::test]
async fn torrent_files_are_multipart() {
    let f = Fake::new().await;
    f.ok(r#"{"files":[{"magnet":"m","id":9}]}"#).await;
    let files = vec![
        TorrentFile { name: "a.torrent".into(), content: b"d4:infoe".to_vec() },
        TorrentFile { name: "b.torrent".into(), content: b"x".to_vec() },
    ];
    f.client.upload_torrent_files(files).await.unwrap();
    let req = f.last().await;
    let body = String::from_utf8_lossy(&req.body);
    assert!(req.headers.get("content-type").unwrap().to_str().unwrap().starts_with("multipart/form-data"));
    assert_eq!(body.matches(r#"name="files[]""#).count(), 2);
    assert!(body.contains(r#"filename="a.torrent""#) && body.contains(r#"filename="b.torrent""#));
}

type Call = Box<dyn Fn(AllDebrid) -> Pin<Box<dyn Future<Output = Result<(), Error>>>>>;
/// An operation, the client call, its expected form parameters and the data to answer.
type Route = (&'static str, Call, &'static [(&'static str, &'static [&'static str])], &'static str);

macro_rules! call {
    ($c:ident => $e:expr) => {
        Box::new(|$c: AllDebrid| Box::pin(async move { $e.await.map(|_| ()) }) as Pin<Box<dyn Future<Output = _>>>)
            as Call
    };
}

#[tokio::test]
async fn routes_every_operation_of_the_spec() {
    let routes: Vec<Route> = vec![
        ("GET /v4/ping", call!(c => c.ping()), &[], "{}"),
        (
            "GET /v4.1/pin/get",
            call!(c => c.get_pin()),
            &[],
            r#"{"pin":"P","check":"C","expires_in":600,"user_url":"u","base_url":"b"}"#,
        ),
        (
            "POST /v4/pin/check",
            call!(c => c.check_pin("P", "C")),
            &[("check", &["C"]), ("pin", &["P"])],
            r#"{"activated":false,"expires_in":600}"#,
        ),
        ("GET /v4/hosts", call!(c => c.hosts(false)), &[], r#"{"hosts":{}}"#),
        ("GET /v4/hosts/domains", call!(c => c.host_domains()), &[], r#"{"hosts":[]}"#),
        ("GET /v4/hosts/priority", call!(c => c.host_priority()), &[], r#"{"hosts":{}}"#),
        ("GET /v4/user", call!(c => c.user()), &[], USER),
        ("GET /v4.1/user/hosts", call!(c => c.user_hosts(false)), &[], r#"{"hosts":{}}"#),
        ("POST /v4/user/verif", call!(c => c.verification_status("T")), &[("token", &["T"])], r#"{"verif":"waiting"}"#),
        ("POST /v4/user/verif/resend", call!(c => c.resend_verification("T")), &[("token", &["T"])], "{}"),
        ("POST /v4/user/notification/clear", call!(c => c.clear_notification("N")), &[("code", &["N"])], "{}"),
        (
            "POST /v4/link/infos",
            call!(c => c.link_infos(&["a", "b"], Some("p"))),
            &[("link[]", &["a", "b"]), ("password", &["p"])],
            r#"{"infos":[]}"#,
        ),
        ("POST /v4/link/redirector", call!(c => c.extract_redirector("L")), &[("link", &["L"])], r#"{"links":[]}"#),
        ("POST /v4/link/unlock", call!(c => c.unlock("L", None)), &[("link", &["L"])], "{}"),
        (
            "POST /v4/link/streaming",
            call!(c => c.streaming_link("G", "S")),
            &[("id", &["G"]), ("stream", &["S"])],
            "{}",
        ),
        (
            "POST /v4/link/delayed",
            call!(c => c.delayed_link("42")),
            &[("id", &["42"])],
            r#"{"status":1,"time_left":5}"#,
        ),
        ("POST /v4/magnet/upload", call!(c => c.upload_magnets(&["M"])), &[("magnets[]", &["M"])], r#"{"magnets":[]}"#),
        (
            "POST /v4/magnet/upload/file",
            call!(c => c.upload_torrent_files(vec![TorrentFile { name: "a".into(), content: vec![1] }])),
            &[],
            r#"{"files":[]}"#,
        ),
        ("POST /v4.1/magnet/status", call!(c => c.magnets()), &[], r#"{"magnets":[]}"#),
        ("POST /v4/magnet/files", call!(c => c.magnet_files(&[1, 2])), &[("id[]", &["1", "2"])], r#"{"magnets":[]}"#),
        ("POST /v4/magnet/delete", call!(c => c.delete_magnet(5)), &[("id", &["5"])], "{}"),
        ("POST /v4/magnet/restart", call!(c => c.restart_magnets(&[5, 6])), &[("ids[]", &["5", "6"])], "{}"),
        ("GET /v4/user/links", call!(c => c.saved_links()), &[], r#"{"links":[]}"#),
        ("POST /v4/user/links/save", call!(c => c.save_links(&["L"])), &[("links[]", &["L"])], "{}"),
        ("POST /v4/user/links/delete", call!(c => c.delete_saved_links(&["L"])), &[("links[]", &["L"])], "{}"),
        ("GET /v4/user/history", call!(c => c.history()), &[], r#"{"links":[]}"#),
        ("POST /v4/user/history/delete", call!(c => c.purge_history()), &[], "{}"),
        ("GET /v4/voucher/balance", call!(c => c.voucher_balance()), &[], r#"{"balance":1.5}"#),
        (
            "POST /v4/voucher/get",
            call!(c => c.vouchers(30, 2)),
            &[("duration", &["30"]), ("nb", &["2"])],
            r#"{"codes":[]}"#,
        ),
        (
            "POST /v4/voucher/generate",
            call!(c => c.generate_vouchers(90, 1)),
            &[("duration", &["90"]), ("nb", &["1"])],
            r#"{"codes":[]}"#,
        ),
    ];
    let ops = spec_operations();
    for (want, call, pairs, data) in &routes {
        let f = Fake::new().await;
        f.ok(data).await;
        call(f.client.clone()).await.unwrap_or_else(|e| panic!("{want}: {e}"));
        let req = f.last().await;
        let got = format!("{} {}", req.method, req.url.path());
        assert_eq!(&got, want);
        assert!(ops.contains(&got), "{got} is not in the spec");
        if !pairs.is_empty() {
            assert_eq!(form(&req), expect_form(pairs), "{want}");
        }
    }
    let covered: BTreeSet<String> = routes.iter().map(|r| r.0.to_owned()).collect();
    assert_eq!(covered, ops);
}

#[test]
fn error_codes_match_the_spec() {
    let text = std::fs::read_to_string(SPEC).unwrap();
    let block = text.split("\n    ErrorCode:\n").nth(1).unwrap();
    let mut codes: Vec<&str> = block
        .lines()
        .take_while(|l| !l.starts_with("    ") || l.starts_with("     "))
        .filter_map(|l| l.strip_prefix("        - "))
        .map(|c| c.trim_matches('"'))
        .collect();
    let mut known = error_codes::ALL.to_vec();
    codes.sort_unstable();
    known.sort_unstable();
    assert_eq!(codes, known);
}

#[tokio::test]
async fn api_errors_carry_their_code() {
    let f = Fake::new().await;
    f.answer(200, r#"{"status":"error","error":{"code":"AUTH_BLOCKED","message":"m","token":"tok"}}"#).await;
    let err = f.client.user().await.unwrap_err();
    assert!(err.is_code(error_codes::AUTH_BLOCKED));
    match err {
        Error::Api { token, .. } => assert_eq!(token.as_deref(), Some("tok")),
        other => panic!("{other:?}"),
    }
}

#[tokio::test]
async fn http_errors() {
    let f = Fake::new().await;
    f.answer(429, "Too Many Requests").await;
    f.answer(404, r#"{"status":"error","error":{"code":"404","message":"Endpoint doesn't exist"}}"#).await;
    f.answer(200, "<html>").await;
    assert!(f.client.ping().await.unwrap_err().is_rate_limited());
    assert!(f.client.ping().await.unwrap_err().is_code(error_codes::NOT_FOUND));
    assert!(matches!(f.client.ping().await.unwrap_err(), Error::Decode(_)));
}

#[tokio::test]
async fn flexible_shapes() {
    let f = Fake::new().await;
    f.ok(r#"{"magnets":{"id":7,"statusCode":4,"files":[{"n":"a","s":1,"l":"L"}]}}"#).await;
    let magnet = f.client.magnet(7).await.unwrap();
    assert_eq!(magnet.id, 7);
    assert!(magnet.is_ready());
    assert_eq!(magnet.files.unwrap()[0].l.as_deref(), Some("L"));

    f.ok(USER).await;
    assert!(f.client.user().await.unwrap().limited_hosters_quotas.is_empty());

    f.ok(r#"{"link":"","delayed":123,"streaming":[{"id":"mp3","quality":"mp3"},{"id":"hd","quality":1080}]}"#).await;
    let unlocked = f.client.unlock("L", None).await.unwrap();
    assert_eq!(unlocked.delayed.unwrap().as_str(), "123");
    let qualities: Vec<String> =
        unlocked.streaming.unwrap().into_iter().map(|s| s.quality.unwrap().to_string()).collect();
    assert_eq!(qualities, ["mp3", "1080"]);

    f.ok(r#"{"hosts":{"a":{"name":"a","type":"free","domains":[],"regexp":"x"}}}"#).await;
    assert_eq!(f.client.hosts(false).await.unwrap().hosts["a"].regexp.as_deref(), Some(&["x".to_owned()][..]));
}

#[tokio::test]
async fn wait_for_pin() {
    let f = Fake::new().await;
    f.ok(r#"{"activated":false,"expires_in":600}"#).await;
    f.ok(r#"{"activated":true,"expires_in":590,"apikey":"K"}"#).await;
    let pin = PinGetData {
        pin: "ABCD".into(),
        check: "chk".into(),
        expires_in: 600,
        user_url: "u".into(),
        base_url: "b".into(),
    };
    assert_eq!(f.client.wait_for_pin(&pin, Duration::from_millis(1)).await.unwrap(), "K");
    assert_eq!(form(&f.last().await), expect_form(&[("check", &["chk"]), ("pin", &["ABCD"])]));
}

#[tokio::test]
async fn wait_delayed_link() {
    let f = Fake::new().await;
    f.ok(r#"{"status":1,"time_left":5}"#).await;
    f.ok(r#"{"status":2,"time_left":0,"link":"L"}"#).await;
    f.ok(r#"{"status":3,"time_left":0}"#).await;
    assert_eq!(f.client.wait_delayed_link("1", Duration::from_millis(1)).await.unwrap(), "L");
    assert!(f.client.wait_delayed_link("1", Duration::from_millis(1)).await.is_err());
}

#[tokio::test]
async fn live_magnets() {
    let f = Fake::new().await;
    f.ok(r#"{"fullsync":true,"counter":1,"magnets":[{"id":1,"filename":"a","uploadDate":10,"downloaded":5},{"id":2,"filename":"b","uploadDate":20}]}"#).await;
    f.ok(r#"{"counter":2,"magnets":{"id":1,"downloaded":9}}"#).await;
    f.ok(r#"{"counter":3,"magnets":[{"id":2,"deleted":true},{"id":3,"filename":"c","uploadDate":30}]}"#).await;
    f.ok(r#"{"counter":4,"magnets":[]}"#).await;
    f.ok(r#"{"fullsync":true,"counter":1,"magnets":[{"id":9,"uploadDate":1}]}"#).await;
    let mut live = f.client.live_magnets();
    let ids = |m: Vec<open_alldebrid::models::Magnet>| m.into_iter().map(|m| m.id).collect::<Vec<_>>();

    assert_eq!(ids(live.poll().await.unwrap()), [2, 1]);
    let first = form(&f.last().await);
    assert_eq!(first["counter"], ["0"]);

    let merged = live.poll().await.unwrap();
    let second = form(&f.last().await);
    assert_eq!(second["counter"], ["1"]);
    assert_eq!(first["session"], second["session"]);
    let one = merged.iter().find(|m| m.id == 1).unwrap();
    assert_eq!((one.filename.as_deref(), one.downloaded), (Some("a"), Some(9)));

    assert_eq!(ids(live.poll().await.unwrap()), [3, 1]);
    assert_eq!(ids(live.poll().await.unwrap()), [3, 1]);
    assert_eq!(ids(live.poll().await.unwrap()), [9]);
}
