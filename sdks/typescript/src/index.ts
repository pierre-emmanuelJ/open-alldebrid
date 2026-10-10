/**
 * Typed client for the AllDebrid API (v4 and v4.1).
 *
 * The types come from `schema.gen.ts`, generated from the open-alldebrid
 * OpenAPI spec; the client is hand-written on top of them, one method per
 * operation of the spec.
 *
 * @module
 */
import { errorCodeValues, type components } from "./schema.gen.js";

type Schemas = components["schemas"];

export type ApiError = Schemas["Error"];
export type ErrorCode = Schemas["ErrorCode"];
export type Service = Schemas["Service"];
export type ServiceMap = Schemas["ServiceMap"];
export type HostsData = Schemas["HostsData"];
export type HostDomainsData = Schemas["HostDomainsData"];
export type PinGetData = Schemas["PinGetData"];
export type PinCheckData = Schemas["PinCheckData"];
export type VerificationData = Schemas["VerificationData"];
export type LinkInfo = Schemas["LinkInfo"];
export type Stream = Schemas["Stream"];
export type UnlockData = Schemas["UnlockData"];
export type StreamingData = Schemas["StreamingData"];
export type DelayedData = Schemas["DelayedData"];
export type UploadedMagnet = Schemas["UploadedMagnet"];
export type Magnet = Schemas["Magnet"];
export type FileNode = Schemas["FileNode"];
export type MagnetFiles = Schemas["MagnetFiles"];
export type RestartedMagnet = Schemas["RestartedMagnet"];
export type SavedLink = Schemas["SavedLink"];
export type VouchersData = Schemas["VouchersData"];
export type GeneratedVouchersData = Schemas["GeneratedVouchersData"];
export type { components } from "./schema.gen.js";

/** The account, with `limitedHostersQuotas` always an object. */
export type User = Omit<Schemas["User"], "limitedHostersQuotas"> & {
  limitedHostersQuotas: Record<string, number>;
};

/** Every error code the spec lists. The API may return others. */
export const ERROR_CODES: readonly ErrorCode[] = errorCodeValues as readonly ErrorCode[];

export const DEFAULT_BASE_URL = "https://api.alldebrid.com";

/** Public demo keys: the API answers them with static sample data. */
export const DEMO_KEYS = {
  premium: "staticDemoApikeyPrem",
  trial: "staticDemoApikeyTria",
  free: "staticDemoApikeyFree",
} as const;

/** Magnet status codes. 5 and above are errors. */
export const MagnetStatus = {
  Queued: 0,
  Downloading: 1,
  Compressing: 2,
  Uploading: 3,
  Ready: 4,
} as const;

/** Delayed link statuses. */
export const DelayedStatus = { Processing: 1, Ready: 2, Failed: 3 } as const;

/** An error answered by the API (`"status": "error"`). */
export class AllDebridError extends Error {
  override readonly name = "AllDebridError";
  /** The error code, e.g. `MAGNET_INVALID_ID`. */
  readonly code: string;
  /** With `AUTH_BLOCKED`, the token to poll with `verificationStatus`. */
  readonly token: string | undefined;

  constructor(error: ApiError) {
    super(`${error.code}: ${error.message}`);
    this.code = error.code;
    this.token = error.token;
  }

  /** Whether this error has the given code. */
  is(code: ErrorCode): boolean {
    return this.code === code;
  }
}

/** An HTTP status other than 200, such as 429 when rate limited. */
export class HttpError extends Error {
  override readonly name = "HttpError";
  readonly status: number;
  /** The API error in the body, when there is one. */
  readonly apiError: AllDebridError | undefined;

  constructor(status: number, apiError?: AllDebridError) {
    super(apiError ? `HTTP ${status}: ${apiError.message}` : `HTTP ${status}`);
    this.status = status;
    this.apiError = apiError;
  }

  /** Whether the API asked to slow down (HTTP 429 or 503). */
  get rateLimited(): boolean {
    return this.status === 429 || this.status === 503;
  }
}

/** Whether `err` is an API error with the given code. */
export function isErrorCode(err: unknown, code: ErrorCode): boolean {
  if (err instanceof HttpError) return err.apiError?.code === code;
  return err instanceof AllDebridError && err.code === code;
}

export interface ClientOptions {
  /** An API key. Leave it out for public endpoints and the PIN flow. */
  apiKey?: string;
  /** Defaults to https://api.alldebrid.com. */
  baseUrl?: string;
  /** Identifies your app. Browsers ignore it. */
  userAgent?: string;
  /** A fetch implementation; defaults to the global fetch. */
  fetch?: typeof fetch;
}

export interface PollOptions {
  /** Delay between two calls, in milliseconds. Defaults to 5000. */
  intervalMs?: number;
  /** Stops polling. */
  signal?: AbortSignal;
}

/** A .torrent file to upload. */
export interface TorrentFile {
  name: string;
  data: Blob | Uint8Array | ArrayBuffer;
}

type Params = [string, string][];

interface Envelope<T> {
  status: "success" | "error";
  data?: T;
  error?: ApiError;
}

/** A client for the AllDebrid API. */
export class AllDebrid {
  readonly #apiKey: string | undefined;
  readonly #baseUrl: string;
  readonly #userAgent: string | undefined;
  readonly #fetch: typeof fetch;

  constructor(options: ClientOptions = {}) {
    this.#apiKey = options.apiKey;
    this.#baseUrl = (options.baseUrl ?? DEFAULT_BASE_URL).replace(/\/+$/, "");
    this.#userAgent = options.userAgent;
    this.#fetch = options.fetch ?? globalThis.fetch.bind(globalThis);
  }

  /** A copy of this client with another API key. */
  withApiKey(apiKey: string): AllDebrid {
    return new AllDebrid({
      apiKey,
      baseUrl: this.#baseUrl,
      ...(this.#userAgent !== undefined && { userAgent: this.#userAgent }),
      fetch: this.#fetch,
    });
  }

  async #call<T>(method: "GET" | "POST", path: string, params: Params = [], body?: FormData): Promise<T> {
    const headers: Record<string, string> = { Accept: "application/json" };
    if (this.#apiKey) headers["Authorization"] = `Bearer ${this.#apiKey}`;
    if (this.#userAgent) headers["User-Agent"] = this.#userAgent;

    let url = this.#baseUrl + path;
    let payload: BodyInit | undefined = body;
    if (!body && params.length > 0) {
      // List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
      const encoded = new URLSearchParams(params);
      if (method === "GET") url += `?${encoded.toString()}`;
      else payload = encoded;
    }

    const response = await this.#fetch(url, { method, headers, ...(payload !== undefined && { body: payload }) });
    const text = await response.text();
    let envelope: Envelope<T> | undefined;
    try {
      envelope = JSON.parse(text) as Envelope<T>;
    } catch {
      envelope = undefined;
    }
    if (response.status !== 200) {
      throw new HttpError(response.status, envelope?.error ? new AllDebridError(envelope.error) : undefined);
    }
    if (!envelope) throw new Error(`AllDebrid ${method} ${path}: the response is not JSON`);
    if (envelope.status !== "success") {
      if (envelope.error) throw new AllDebridError(envelope.error);
      throw new Error(`AllDebrid ${method} ${path}: status ${String(envelope.status)} without error`);
    }
    return envelope.data as T;
  }

  // ------------------------------------------------------------------ general

  /** Checks that the API is up. */
  async ping(): Promise<void> {
    await this.#call("GET", "/v4/ping");
  }

  // --------------------------------------------------------------------- pin

  /** Starts a PIN authentication. Show `pin` or `user_url` to the user, then `waitForPin`. */
  getPin(): Promise<PinGetData> {
    return this.#call("GET", "/v4.1/pin/get");
  }

  /** The state of a PIN authentication; `apikey` is set once activated. */
  checkPin(pin: string, check: string): Promise<PinCheckData> {
    return this.#call("POST", "/v4/pin/check", [["pin", pin], ["check", check]]);
  }

  /** Polls until the user confirms the PIN and returns the new API key. */
  async waitForPin(pin: Pick<PinGetData, "pin" | "check">, options: PollOptions = {}): Promise<string> {
    for (;;) {
      const state = await this.checkPin(pin.pin, pin.check);
      if (state.activated && state.apikey) return state.apikey;
      await sleep(options.intervalMs ?? 5000, options.signal);
    }
  }

  // ------------------------------------------------------------------- hosts

  /** Supported hosts, streaming sites and redirectors. Cache the result. */
  hosts(options: { hostsOnly?: boolean } = {}): Promise<HostsData> {
    return this.#call("GET", "/v4/hosts", options.hostsOnly ? [["hostsOnly", "true"]] : []);
  }

  /** Every supported domain. */
  hostDomains(): Promise<HostDomainsData> {
    return this.#call("GET", "/v4/hosts/domains");
  }

  /** Host names ranked from the most open (1) to the most restricted. */
  async hostPriority(): Promise<Record<string, number>> {
    return (await this.#call<Schemas["HostPriorityData"]>("GET", "/v4/hosts/priority")).hosts;
  }

  // -------------------------------------------------------------------- user

  /** The account information. */
  async user(): Promise<User> {
    const { user } = await this.#call<Schemas["UserData"]>("GET", "/v4/user");
    const quotas = Array.isArray(user.limitedHostersQuotas) ? {} : user.limitedHostersQuotas;
    return { ...user, limitedHostersQuotas: quotas };
  }

  /** The hosts this user can use now, with their quotas. */
  async userHosts(options: { hostsOnly?: boolean } = {}): Promise<ServiceMap> {
    const params: Params = options.hostsOnly ? [["hostsOnly", "true"]] : [];
    return (await this.#call<Schemas["UserHostsData"]>("GET", "/v4.1/user/hosts", params)).hosts;
  }

  /** Follows a new-location login verification (token from an AUTH_BLOCKED error). */
  verificationStatus(token: string): Promise<VerificationData> {
    return this.#call("POST", "/v4/user/verif", [["token", token]]);
  }

  /** Sends the verification email again (once per token). */
  async resendVerification(token: string): Promise<void> {
    await this.#call("POST", "/v4/user/verif/resend", [["token", token]]);
  }

  /** Clears one of the codes in `user().notifications`. */
  async clearNotification(code: string): Promise<void> {
    await this.#call("POST", "/v4/user/notification/clear", [["code", code]]);
  }

  // ------------------------------------------------------------------- links

  /** File name and size of each link, in order; unreadable links have `error` set. */
  async linkInfos(links: string[], options: { password?: string } = {}): Promise<LinkInfo[]> {
    const params: Params = links.map((l) => ["link[]", l]);
    if (options.password) params.push(["password", options.password]);
    return (await this.#call<Schemas["LinkInfosData"]>("POST", "/v4/link/infos", params)).infos;
  }

  /** Links behind a redirector; they only work with `linkInfos` and `unlock`. */
  async extractRedirector(link: string): Promise<string[]> {
    return (await this.#call<Schemas["RedirectorData"]>("POST", "/v4/link/redirector", [["link", link]])).links;
  }

  /**
   * Unlocks a hoster link. When `delayed` is set, wait with `waitDelayedLink`;
   * when `streaming` lists qualities, choose one with `streamingLink`.
   */
  unlock(link: string, options: { password?: string } = {}): Promise<UnlockData> {
    const params: Params = [["link", link]];
    if (options.password) params.push(["password", options.password]);
    return this.#call("POST", "/v4/link/unlock", params);
  }

  /** The link of one quality: `id` from `unlock`, `stream` from one of its `streaming` entries. */
  streamingLink(id: string, stream: string): Promise<StreamingData> {
    return this.#call("POST", "/v4/link/streaming", [["id", id], ["stream", stream]]);
  }

  /** The status of a delayed link. */
  delayedLink(id: string | number): Promise<DelayedData> {
    return this.#call("POST", "/v4/link/delayed", [["id", String(id)]]);
  }

  /** Polls a delayed link until it is ready and returns the download link. */
  async waitDelayedLink(id: string | number, options: PollOptions = {}): Promise<string> {
    for (;;) {
      const state = await this.delayedLink(id);
      if (state.status === DelayedStatus.Ready && state.link) return state.link;
      if (state.status === DelayedStatus.Failed) throw new Error(`AllDebrid: delayed link ${id} failed`);
      await sleep(options.intervalMs ?? 5000, options.signal);
    }
  }

  // ----------------------------------------------------------------- magnets

  /** Adds magnet URIs or info hashes. Results are in order; invalid ones have `error` set. */
  async uploadMagnets(magnets: string[]): Promise<UploadedMagnet[]> {
    const params: Params = magnets.map((m) => ["magnets[]", m]);
    return (await this.#call<Schemas["UploadMagnetsData"]>("POST", "/v4/magnet/upload", params)).magnets;
  }

  /** Adds .torrent files. Results are in order; invalid files have `error` set. */
  async uploadTorrentFiles(files: TorrentFile[]): Promise<UploadedMagnet[]> {
    const form = new FormData();
    for (const f of files) {
      const blob = f.data instanceof Blob ? f.data : new Blob([f.data as BlobPart], { type: "application/x-bittorrent" });
      form.append("files[]", blob, f.name);
    }
    return (await this.#call<Schemas["UploadTorrentFilesData"]>("POST", "/v4/magnet/upload/file", [], form)).files;
  }

  /** Every magnet of the account. */
  async magnets(): Promise<Magnet[]> {
    return toList((await this.#call<Schemas["MagnetStatusData"]>("POST", "/v4.1/magnet/status")).magnets);
  }

  /** One magnet, with its `files`. */
  async magnet(id: number): Promise<Magnet> {
    const data = await this.#call<Schemas["MagnetStatusData"]>("POST", "/v4.1/magnet/status", [["id", String(id)]]);
    const [magnet, ...rest] = toList(data.magnets);
    if (!magnet || rest.length > 0) throw new Error(`AllDebrid: magnet ${id}: unexpected response`);
    return magnet;
  }

  /** The file tree of each magnet, in order; unknown IDs have `error` set. */
  async magnetFiles(ids: number[]): Promise<MagnetFiles[]> {
    const params: Params = ids.map((id) => ["id[]", String(id)]);
    return (await this.#call<Schemas["MagnetFilesData"]>("POST", "/v4/magnet/files", params)).magnets;
  }

  /** Deletes a magnet. */
  async deleteMagnet(id: number): Promise<void> {
    await this.#call("POST", "/v4/magnet/delete", [["id", String(id)]]);
  }

  /** Restarts a failed magnet. */
  async restartMagnet(id: number): Promise<void> {
    await this.#call("POST", "/v4/magnet/restart", [["id", String(id)]]);
  }

  /** Restarts failed magnets; each result has its own `error`. */
  async restartMagnets(ids: number[]): Promise<RestartedMagnet[]> {
    const params: Params = ids.map((id) => ["ids[]", String(id)]);
    return (await this.#call<Schemas["RestartMagnetsData"]>("POST", "/v4/magnet/restart", params)).magnets ?? [];
  }

  /** Follows the magnets with live mode: each `poll` only transfers changes. */
  liveMagnets(): LiveMagnets {
    return new LiveMagnets((session, counter) =>
      this.#call<LiveData>("POST", "/v4.1/magnet/status", [
        ["session", String(session)],
        ["counter", String(counter)],
      ]),
    );
  }

  // ------------------------------------------------------------- saved links

  /** The links the user saved. */
  async savedLinks(): Promise<SavedLink[]> {
    return (await this.#call<Schemas["SavedLinksData"]>("GET", "/v4/user/links")).links;
  }

  /** Saves links for later. */
  async saveLinks(links: string[]): Promise<void> {
    await this.#call("POST", "/v4/user/links/save", links.map((l) => ["links[]", l]));
  }

  /** Deletes saved links. */
  async deleteSavedLinks(links: string[]): Promise<void> {
    await this.#call("POST", "/v4/user/links/delete", links.map((l) => ["links[]", l]));
  }

  /** Links unlocked in the last 3 days, when history logging is enabled. */
  async history(): Promise<SavedLink[]> {
    return (await this.#call<Schemas["SavedLinksData"]>("GET", "/v4/user/history")).links;
  }

  /** Deletes the whole recent links history. */
  async purgeHistory(): Promise<void> {
    await this.#call("POST", "/v4/user/history/delete");
  }

  // ---------------------------------------------------------------- reseller

  /** The reseller balance. Reseller accounts only. */
  async voucherBalance(): Promise<number> {
    return (await this.#call<Schemas["VoucherBalanceData"]>("GET", "/v4/voucher/balance")).balance;
  }

  /** Takes vouchers from the reseller stock. Reseller accounts only. */
  vouchers(durationDays: 15 | 30 | 90 | 180 | 365, count: number): Promise<VouchersData> {
    return this.#call("POST", "/v4/voucher/get", [["duration", String(durationDays)], ["nb", String(count)]]);
  }

  /** Generates new vouchers from the reseller balance. Reseller accounts only. */
  generateVouchers(durationDays: 15 | 30 | 90 | 180 | 365, count: number): Promise<GeneratedVouchersData> {
    return this.#call("POST", "/v4/voucher/generate", [["duration", String(durationDays)], ["nb", String(count)]]);
  }
}

interface LiveData {
  magnets?: Magnet[] | Magnet | null;
  counter: number;
  fullsync?: boolean;
}

/** Live mode state: the merged list of magnets, updated by `poll`. */
export class LiveMagnets {
  readonly session = Math.floor(Math.random() * 2 ** 31);
  #counter = 0;
  #state = new Map<number, Magnet>();
  readonly #fetch: (session: number, counter: number) => Promise<LiveData>;

  constructor(fetchStatus: (session: number, counter: number) => Promise<LiveData>) {
    this.#fetch = fetchStatus;
  }

  /** Applies the changes since the last call and returns every magnet, newest first. */
  async poll(): Promise<Magnet[]> {
    const data = await this.#fetch(this.session, this.#counter);
    if (data.fullsync) this.#state = new Map();
    for (const change of toList(data.magnets)) {
      if (change.deleted) this.#state.delete(change.id);
      else this.#state.set(change.id, { ...this.#state.get(change.id), ...change });
    }
    this.#counter = data.counter;
    return [...this.#state.values()].sort((a, b) => (b.uploadDate ?? 0) - (a.uploadDate ?? 0) || b.id - a.id);
  }
}

function toList<T>(value: T[] | T | null | undefined): T[] {
  if (value === null || value === undefined) return [];
  return Array.isArray(value) ? value : [value];
}

function sleep(ms: number, signal?: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) return reject(signal.reason);
    const timer = setTimeout(() => {
      signal?.removeEventListener("abort", onAbort);
      resolve();
    }, ms);
    const onAbort = () => {
      clearTimeout(timer);
      reject(signal?.reason);
    };
    signal?.addEventListener("abort", onAbort, { once: true });
  });
}
