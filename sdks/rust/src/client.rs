//! The API client.

use std::collections::HashMap;
use std::time::Duration;

use reqwest::multipart;
use serde::de::DeserializeOwned;
use serde_json::Value;

use crate::error::{Error, Result};
use crate::live::LiveMagnets;
use crate::models::{
    DelayedData, GeneratedVouchersData, HostDomainsData, HostPriorityData, HostsData, LinkInfo, LinkInfosData, Magnet,
    MagnetFiles, MagnetFilesData, MagnetStatusData, PinCheckData, PinGetData, RedirectorData, RestartMagnetsData,
    RestartedMagnet, SavedLink, SavedLinksData, ServiceMap, StreamingData, UnlockData, UploadMagnetsData,
    UploadTorrentFilesData, UploadedMagnet, User, UserData, UserHostsData, VerificationData, VoucherBalanceData,
    VouchersData,
};

/// The production API.
pub const DEFAULT_BASE_URL: &str = "https://api.alldebrid.com";
const DEFAULT_USER_AGENT: &str = "open-alldebrid-rust";

/// Public demo keys: the API answers them with static sample data.
pub mod demo_keys {
    /// A premium account.
    pub const PREMIUM: &str = "staticDemoApikeyPrem";
    /// An account in its free trial.
    pub const TRIAL: &str = "staticDemoApikeyTria";
    /// A free account.
    pub const FREE: &str = "staticDemoApikeyFree";
}

/// Magnet status codes. 5 and above are errors.
pub mod magnet_status {
    /// In the queue.
    pub const QUEUED: i64 = 0;
    /// Downloading.
    pub const DOWNLOADING: i64 = 1;
    /// Compressing or moving.
    pub const COMPRESSING: i64 = 2;
    /// Uploading.
    pub const UPLOADING: i64 = 3;
    /// The files are available.
    pub const READY: i64 = 4;
}

/// Delayed link statuses.
pub mod delayed_status {
    /// Still processing.
    pub const PROCESSING: i64 = 1;
    /// The link is ready.
    pub const READY: i64 = 2;
    /// The link could not be generated.
    pub const FAILED: i64 = 3;
}

impl Magnet {
    /// Whether the magnet's files are available.
    #[must_use]
    pub fn is_ready(&self) -> bool {
        self.status_code == Some(magnet_status::READY)
    }

    /// Whether the magnet failed.
    #[must_use]
    pub fn is_error(&self) -> bool {
        self.status_code.unwrap_or(0) > magnet_status::READY
    }
}

/// A .torrent file to upload.
#[derive(Debug, Clone)]
pub struct TorrentFile {
    /// The file name.
    pub name: String,
    /// The content of the file.
    pub content: Vec<u8>,
}

/// A client for the AllDebrid API.
///
/// ```no_run
/// # async fn run() -> open_alldebrid::Result<()> {
/// let client = open_alldebrid::AllDebrid::new("api key").user_agent("my-app");
/// let user = client.user().await?;
/// # Ok(()) }
/// ```
///
/// API errors are [`Error::Api`]; other HTTP statuses [`Error::Http`].
#[derive(Debug, Clone)]
pub struct AllDebrid {
    api_key: Option<String>,
    base_url: String,
    user_agent: String,
    http: reqwest::Client,
}

type Params = Vec<(String, String)>;

fn param<T: ToString + ?Sized>(name: &str, value: &T) -> (String, String) {
    (name.to_owned(), value.to_string())
}

fn repeat<T: ToString>(name: &str, values: impl IntoIterator<Item = T>) -> Params {
    values.into_iter().map(|v| param(name, &v)).collect()
}

impl AllDebrid {
    /// A client for the given API key.
    pub fn new(api_key: impl Into<String>) -> Self {
        Self::with_optional_key(Some(api_key.into()))
    }

    /// A client without a key, for public endpoints and the PIN flow.
    #[must_use]
    pub fn anonymous() -> Self {
        Self::with_optional_key(None)
    }

    fn with_optional_key(api_key: Option<String>) -> Self {
        Self {
            api_key: api_key.filter(|k| !k.is_empty()),
            base_url: DEFAULT_BASE_URL.to_owned(),
            user_agent: DEFAULT_USER_AGENT.to_owned(),
            http: reqwest::Client::new(),
        }
    }

    /// Identifies your application to AllDebrid.
    #[must_use]
    pub fn user_agent(mut self, user_agent: impl Into<String>) -> Self {
        self.user_agent = user_agent.into();
        self
    }

    /// Sets the API base URL, for tests or a proxy.
    #[must_use]
    pub fn base_url(mut self, base_url: impl Into<String>) -> Self {
        base_url.into().trim_end_matches('/').clone_into(&mut self.base_url);
        self
    }

    /// Uses this HTTP client, for timeouts or a proxy.
    #[must_use]
    pub fn http_client(mut self, http: reqwest::Client) -> Self {
        self.http = http;
        self
    }

    /// A client with another API key, sharing this one's HTTP client.
    #[must_use]
    pub fn with_api_key(&self, api_key: Option<&str>) -> Self {
        Self { api_key: api_key.filter(|k| !k.is_empty()).map(ToOwned::to_owned), ..self.clone() }
    }

    /// Calls the API and returns the `data` of a successful answer.
    pub(crate) async fn call_raw(
        &self,
        method: reqwest::Method,
        path: &str,
        params: Params,
        form: Option<multipart::Form>,
    ) -> Result<Value> {
        let url = format!("{}{}", self.base_url, path);
        let mut request = self
            .http
            .request(method.clone(), url)
            .header("Accept", "application/json")
            .header("User-Agent", &self.user_agent);
        if let Some(key) = &self.api_key {
            request = request.bearer_auth(key);
        }
        // List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
        request = match form {
            Some(form) => request.multipart(form),
            None if method == reqwest::Method::GET => request.query(&params),
            None => request.form(&params),
        };

        let response = request.send().await?;
        let status = response.status().as_u16();
        let text = response.text().await?;
        let envelope: Option<Value> = serde_json::from_str(&text).ok();
        let api_error = envelope.as_ref().and_then(|e| e.get("error")).and_then(to_error);
        if status != 200 {
            return Err(Error::Http { status, code: api_error.and_then(|e| e.code().map(ToOwned::to_owned)) });
        }
        let Some(Value::Object(mut envelope)) = envelope else {
            return Err(Error::Decode(format!("{path}: the response is not a JSON object")));
        };
        if envelope.get("status").and_then(Value::as_str) != Some("success") {
            return Err(api_error.unwrap_or_else(|| Error::Decode(format!("{path}: error without details"))));
        }
        Ok(envelope.remove("data").unwrap_or(Value::Null))
    }

    async fn call<T: DeserializeOwned>(&self, method: reqwest::Method, path: &str, params: Params) -> Result<T> {
        let data = self.call_raw(method, path, params, None).await?;
        serde_json::from_value(data).map_err(|e| Error::Decode(format!("{path}: {e}")))
    }

    async fn get<T: DeserializeOwned>(&self, path: &str, params: Params) -> Result<T> {
        self.call(reqwest::Method::GET, path, params).await
    }

    async fn post<T: DeserializeOwned>(&self, path: &str, params: Params) -> Result<T> {
        self.call(reqwest::Method::POST, path, params).await
    }

    async fn post_unit(&self, path: &str, params: Params) -> Result<()> {
        self.call_raw(reqwest::Method::POST, path, params, None).await.map(|_| ())
    }

    // ------------------------------------------------------------------ general

    /// Checks that the API is up.
    pub async fn ping(&self) -> Result<()> {
        self.call_raw(reqwest::Method::GET, "/v4/ping", Vec::new(), None).await.map(|_| ())
    }

    // ---------------------------------------------------------------------- pin

    /// Starts a PIN authentication. Show `pin` or `user_url`, then [`Self::wait_for_pin`].
    pub async fn get_pin(&self) -> Result<PinGetData> {
        self.get("/v4.1/pin/get", Vec::new()).await
    }

    /// The state of a PIN authentication; `apikey` is set once activated.
    pub async fn check_pin(&self, pin: &str, check: &str) -> Result<PinCheckData> {
        self.post("/v4/pin/check", vec![param("pin", pin), param("check", check)]).await
    }

    /// Polls until the user confirms the PIN and returns the new API key.
    pub async fn wait_for_pin(&self, pin: &PinGetData, interval: Duration) -> Result<String> {
        loop {
            let state = self.check_pin(&pin.pin, &pin.check).await?;
            if let (true, Some(key)) = (state.activated, state.apikey) {
                if !key.is_empty() {
                    return Ok(key);
                }
            }
            tokio::time::sleep(interval).await;
        }
    }

    // -------------------------------------------------------------------- hosts

    /// Supported hosts, streaming sites and redirectors. Cache the result.
    pub async fn hosts(&self, hosts_only: bool) -> Result<HostsData> {
        self.get("/v4/hosts", hosts_only_param(hosts_only)).await
    }

    /// Every supported domain.
    pub async fn host_domains(&self) -> Result<HostDomainsData> {
        self.get("/v4/hosts/domains", Vec::new()).await
    }

    /// Host names ranked from the most open (1) to the most restricted.
    pub async fn host_priority(&self) -> Result<HashMap<String, i64>> {
        Ok(self.get::<HostPriorityData>("/v4/hosts/priority", Vec::new()).await?.hosts)
    }

    // --------------------------------------------------------------------- user

    /// The account information.
    pub async fn user(&self) -> Result<User> {
        Ok(self.get::<UserData>("/v4/user", Vec::new()).await?.user)
    }

    /// The hosts this user can use now, with their quotas.
    pub async fn user_hosts(&self, hosts_only: bool) -> Result<ServiceMap> {
        Ok(self.get::<UserHostsData>("/v4.1/user/hosts", hosts_only_param(hosts_only)).await?.hosts)
    }

    /// Follows a new-location login verification (token from an `AUTH_BLOCKED` error).
    pub async fn verification_status(&self, token: &str) -> Result<VerificationData> {
        self.post("/v4/user/verif", vec![param("token", token)]).await
    }

    /// Sends the verification email again (once per token).
    pub async fn resend_verification(&self, token: &str) -> Result<()> {
        self.post_unit("/v4/user/verif/resend", vec![param("token", token)]).await
    }

    /// Clears one of the codes in `User::notifications`.
    pub async fn clear_notification(&self, code: &str) -> Result<()> {
        self.post_unit("/v4/user/notification/clear", vec![param("code", code)]).await
    }

    // -------------------------------------------------------------------- links

    /// File name and size of each link, in order; unreadable links have an `error`.
    pub async fn link_infos(&self, links: &[&str], password: Option<&str>) -> Result<Vec<LinkInfo>> {
        let mut params = repeat("link[]", links);
        if let Some(password) = password {
            params.push(param("password", password));
        }
        Ok(self.post::<LinkInfosData>("/v4/link/infos", params).await?.infos)
    }

    /// Links behind a redirector; they only work with [`Self::link_infos`] and [`Self::unlock`].
    pub async fn extract_redirector(&self, link: &str) -> Result<Vec<String>> {
        Ok(self.post::<RedirectorData>("/v4/link/redirector", vec![param("link", link)]).await?.links)
    }

    /// Unlocks a hoster link. When `delayed` is set, wait with
    /// [`Self::wait_delayed_link`]; when `streaming` lists qualities, choose one
    /// with [`Self::streaming_link`].
    pub async fn unlock(&self, link: &str, password: Option<&str>) -> Result<UnlockData> {
        let mut params = vec![param("link", link)];
        if let Some(password) = password {
            params.push(param("password", password));
        }
        self.post("/v4/link/unlock", params).await
    }

    /// The link of one quality: `id` from [`Self::unlock`], `stream` from one of its streaming entries.
    pub async fn streaming_link(&self, id: &str, stream: &str) -> Result<StreamingData> {
        self.post("/v4/link/streaming", vec![param("id", &id), param("stream", stream)]).await
    }

    /// The status of a delayed link.
    pub async fn delayed_link(&self, id: &str) -> Result<DelayedData> {
        self.post("/v4/link/delayed", vec![param("id", &id)]).await
    }

    /// Polls a delayed link until it is ready and returns the download link.
    pub async fn wait_delayed_link(&self, id: &str, interval: Duration) -> Result<String> {
        loop {
            let state = self.delayed_link(id).await?;
            match (state.status, state.link) {
                (delayed_status::READY, Some(link)) if !link.is_empty() => return Ok(link),
                (delayed_status::FAILED, _) => return Err(Error::Decode(format!("delayed link {id} failed"))),
                _ => tokio::time::sleep(interval).await,
            }
        }
    }

    // ------------------------------------------------------------------ magnets

    /// Adds magnet URIs or info hashes. Results are in order; invalid ones have an `error`.
    pub async fn upload_magnets(&self, magnets: &[&str]) -> Result<Vec<UploadedMagnet>> {
        Ok(self.post::<UploadMagnetsData>("/v4/magnet/upload", repeat("magnets[]", magnets)).await?.magnets)
    }

    /// Adds .torrent files. Results are in order; invalid files have an `error`.
    pub async fn upload_torrent_files(&self, files: Vec<TorrentFile>) -> Result<Vec<UploadedMagnet>> {
        let mut form = multipart::Form::new();
        for file in files {
            let part =
                multipart::Part::bytes(file.content).file_name(file.name).mime_str("application/x-bittorrent")?;
            form = form.part("files[]", part);
        }
        let data = self.call_raw(reqwest::Method::POST, "/v4/magnet/upload/file", Vec::new(), Some(form)).await?;
        let data: UploadTorrentFilesData = serde_json::from_value(data).map_err(|e| Error::Decode(e.to_string()))?;
        Ok(data.files)
    }

    /// Every magnet of the account.
    pub async fn magnets(&self) -> Result<Vec<Magnet>> {
        Ok(self.post::<MagnetStatusData>("/v4.1/magnet/status", Vec::new()).await?.magnets)
    }

    /// One magnet, with its files.
    pub async fn magnet(&self, id: i64) -> Result<Magnet> {
        let mut found = self.post::<MagnetStatusData>("/v4.1/magnet/status", vec![param("id", &id)]).await?.magnets;
        if found.len() != 1 {
            return Err(Error::Decode(format!("magnet {id}: got {} magnets", found.len())));
        }
        Ok(found.remove(0))
    }

    /// The file tree of each magnet, in order; unknown IDs have an `error`.
    pub async fn magnet_files(&self, ids: &[i64]) -> Result<Vec<MagnetFiles>> {
        Ok(self.post::<MagnetFilesData>("/v4/magnet/files", repeat("id[]", ids)).await?.magnets)
    }

    /// Deletes a magnet.
    pub async fn delete_magnet(&self, id: i64) -> Result<()> {
        self.post_unit("/v4/magnet/delete", vec![param("id", &id)]).await
    }

    /// Restarts a failed magnet.
    pub async fn restart_magnet(&self, id: i64) -> Result<()> {
        self.post_unit("/v4/magnet/restart", vec![param("id", &id)]).await
    }

    /// Restarts failed magnets; each result has its own `error`.
    pub async fn restart_magnets(&self, ids: &[i64]) -> Result<Vec<RestartedMagnet>> {
        let data = self.post::<RestartMagnetsData>("/v4/magnet/restart", repeat("ids[]", ids)).await?;
        Ok(data.magnets.unwrap_or_default())
    }

    /// Follows the magnets with live mode: each [`LiveMagnets::poll`] only transfers changes.
    #[must_use]
    pub fn live_magnets(&self) -> LiveMagnets {
        LiveMagnets::new(self.clone())
    }

    // -------------------------------------------------------------- saved links

    /// The links the user saved.
    pub async fn saved_links(&self) -> Result<Vec<SavedLink>> {
        Ok(self.get::<SavedLinksData>("/v4/user/links", Vec::new()).await?.links)
    }

    /// Saves links for later.
    pub async fn save_links(&self, links: &[&str]) -> Result<()> {
        self.post_unit("/v4/user/links/save", repeat("links[]", links)).await
    }

    /// Deletes saved links.
    pub async fn delete_saved_links(&self, links: &[&str]) -> Result<()> {
        self.post_unit("/v4/user/links/delete", repeat("links[]", links)).await
    }

    /// Links unlocked in the last 3 days, when history logging is enabled.
    pub async fn history(&self) -> Result<Vec<SavedLink>> {
        Ok(self.get::<SavedLinksData>("/v4/user/history", Vec::new()).await?.links)
    }

    /// Deletes the whole recent links history.
    pub async fn purge_history(&self) -> Result<()> {
        self.post_unit("/v4/user/history/delete", Vec::new()).await
    }

    // ----------------------------------------------------------------- reseller

    /// The reseller balance. Reseller accounts only.
    pub async fn voucher_balance(&self) -> Result<f64> {
        Ok(self.get::<VoucherBalanceData>("/v4/voucher/balance", Vec::new()).await?.balance)
    }

    /// Takes vouchers (15, 30, 90, 180 or 365 days) from the reseller stock.
    pub async fn vouchers(&self, duration_days: u32, count: u32) -> Result<VouchersData> {
        self.post("/v4/voucher/get", voucher_params(duration_days, count)).await
    }

    /// Generates new vouchers from the reseller balance.
    pub async fn generate_vouchers(&self, duration_days: u32, count: u32) -> Result<GeneratedVouchersData> {
        self.post("/v4/voucher/generate", voucher_params(duration_days, count)).await
    }
}

fn hosts_only_param(hosts_only: bool) -> Params {
    if hosts_only { vec![param("hostsOnly", "true")] } else { Vec::new() }
}

fn voucher_params(duration_days: u32, count: u32) -> Params {
    vec![param("duration", &duration_days), param("nb", &count)]
}

fn to_error(error: &Value) -> Option<Error> {
    let error = error.as_object()?;
    let text = |key: &str| error.get(key).and_then(Value::as_str).map(ToOwned::to_owned);
    Some(Error::Api {
        code: text("code").unwrap_or_default(),
        message: text("message").unwrap_or_default(),
        token: text("token"),
    })
}
