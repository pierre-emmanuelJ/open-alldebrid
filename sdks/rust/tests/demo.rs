//! Calls the real API with the public demo keys. Set OFFLINE=1 to skip.

use open_alldebrid::{AllDebrid, TorrentFile, delayed_status, demo_keys, error_codes};

fn offline() -> bool {
    std::env::var("OFFLINE").as_deref() == Ok("1")
}

fn client() -> AllDebrid {
    AllDebrid::new(demo_keys::PREMIUM).user_agent("open-alldebrid-rust-tests")
}

fn anonymous() -> AllDebrid {
    AllDebrid::anonymous().user_agent("open-alldebrid-rust-tests")
}

#[tokio::test]
async fn public_endpoints() {
    if offline() {
        return;
    }
    let c = anonymous();
    c.ping().await.unwrap();
    assert!(!c.hosts(false).await.unwrap().hosts.is_empty());
    assert!(!c.host_domains().await.unwrap().hosts.is_empty());
    assert!(!c.host_priority().await.unwrap().is_empty());
}

#[tokio::test]
async fn pin() {
    if offline() {
        return;
    }
    let c = anonymous();
    assert!(c.get_pin().await.unwrap().user_url.starts_with("https://"));
    assert!(c.check_pin("123B", "664c3ca2635c99f291d28e11ea18e154750bd21a").await.unwrap().activated);
}

#[tokio::test]
async fn user() {
    if offline() {
        return;
    }
    let c = client();
    assert_eq!(c.user().await.unwrap().username, "demoUserPremium");
    assert!(!c.user_hosts(false).await.unwrap().is_empty());
    c.clear_notification("NOTIF_CODE").await.unwrap();
    assert_eq!(anonymous().verification_status("verificationTokenAllowed3f9eb0db").await.unwrap().verif, "allowed");
    let err = anonymous().resend_verification("verificationTokenAlreadySentb0db").await.unwrap_err();
    assert!(err.is_code(error_codes::ALREADY_SENT));
}

#[tokio::test]
async fn links() {
    if offline() {
        return;
    }
    let c = client();
    let infos = c.link_infos(&["https://example.com/file2MB.txt", "https://example.com/down"], None).await.unwrap();
    assert!(infos[0].filename.is_some());
    assert!(!c.extract_redirector("https://example.com/redir/folder").await.unwrap().is_empty());
    assert!(!c.unlock("https://example.com/streaming", None).await.unwrap().streaming.unwrap().is_empty());
    assert!(c.unlock("https://example.com/delayed", None).await.unwrap().delayed.is_some());
    assert!(c.unlock("https://example.com/down", None).await.unwrap_err().is_code(error_codes::LINK_DOWN));
    assert!(c.streaming_link("demov1bc797", "hd-720").await.unwrap().link.is_some());
    assert_eq!(c.delayed_link("1234567").await.unwrap().status, delayed_status::READY);
}

#[tokio::test]
async fn magnets() {
    if offline() {
        return;
    }
    let c = client();
    let uploaded = c.upload_magnets(&["magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707"]).await.unwrap();
    assert!(uploaded[0].id.is_some());
    let file = TorrentFile { name: "demo.torrent".into(), content: b"d4:infod4:name4:demoee".to_vec() };
    c.upload_torrent_files(vec![file]).await.unwrap();
    assert!(!c.magnets().await.unwrap().is_empty());
    assert!(!c.magnet(56789).await.unwrap().files.unwrap().is_empty());
    assert_eq!(c.magnet_files(&[56789]).await.unwrap().len(), 1);
    c.delete_magnet(56789).await.unwrap();
    c.restart_magnet(99999).await.unwrap();
    assert!(!c.live_magnets().poll().await.unwrap().is_empty());
}

#[tokio::test]
async fn saved_links() {
    if offline() {
        return;
    }
    let c = client();
    assert!(!c.saved_links().await.unwrap().is_empty());
    c.save_links(&["https://example.com/file2MB.txt"]).await.unwrap();
    c.delete_saved_links(&["https://example.com/file2MB.txt"]).await.unwrap();
    c.history().await.unwrap();
    c.purge_history().await.unwrap();
}

#[tokio::test]
async fn bad_key() {
    if offline() {
        return;
    }
    let err = client().with_api_key(Some("not-a-real-key")).user().await.unwrap_err();
    assert!(err.is_code(error_codes::AUTH_BAD_APIKEY));
}
