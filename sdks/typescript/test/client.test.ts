import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { parse } from "yaml";
import { describe, expect, it } from "vitest";
import {
  AllDebrid,
  AllDebridError,
  ERROR_CODES,
  HttpError,
  isErrorCode,
  MagnetStatus,
} from "../src/index.js";

const spec = parse(readFileSync(fileURLToPath(new URL("../../../openapi.yaml", import.meta.url)), "utf8")) as {
  paths: Record<string, Record<string, { operationId: string }>>;
  components: { schemas: { ErrorCode: { enum: string[] } } };
};

interface Recorded {
  method: string;
  path: string;
  query: URLSearchParams;
  form: URLSearchParams;
  files: [string, string][];
  headers: Headers;
}

/** A fetch that records requests and answers canned bodies in order. */
function fake(...answers: { status?: number; body: string }[]) {
  const requests: Recorded[] = [];
  const fetchFn = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const url = new URL(String(input));
    const headers = new Headers(init?.headers);
    let form = new URLSearchParams();
    const files: [string, string][] = [];
    if (init?.body instanceof URLSearchParams) form = init.body;
    if (init?.body instanceof FormData) {
      for (const [k, v] of init.body.entries()) if (v instanceof Blob) files.push([k, (v as File).name]);
    }
    requests.push({ method: init?.method ?? "GET", path: url.pathname, query: url.searchParams, form, files, headers });
    const a = answers.shift() ?? { body: '{"status":"success","data":{}}' };
    return new Response(a.body, { status: a.status ?? 200, headers: { "Content-Type": "application/json" } });
  };
  const client = new AllDebrid({ apiKey: "test-key", baseUrl: "https://api.test", userAgent: "unit", fetch: fetchFn });
  return { client, requests, last: () => requests[requests.length - 1]! };
}

const ok = (data: unknown) => ({ body: JSON.stringify({ status: "success", data }) });
const apiError = (code: string, extra: object = {}) => ({
  body: JSON.stringify({ status: "error", error: { code, message: "m", ...extra } }),
});

describe("requests", () => {
  it("sends list parameters PHP-style, with the key in a header", async () => {
    const { client, last } = fake(ok({ magnets: [{ magnet: "a", id: 1 }] }));
    const result = await client.uploadMagnets(["a", "b"]);
    expect(last().method).toBe("POST");
    expect(last().path).toBe("/v4/magnet/upload");
    expect(last().form.getAll("magnets[]")).toEqual(["a", "b"]);
    expect(last().headers.get("Authorization")).toBe("Bearer test-key");
    expect(last().headers.get("User-Agent")).toBe("unit");
    expect(result[0]?.id).toBe(1);
  });

  it("omits the Authorization header without a key", async () => {
    const { client, last } = fake(ok({ ping: "pong" }));
    await client.withApiKey("").ping();
    expect(last().headers.get("Authorization")).toBeNull();
  });

  it("sends GET parameters in the query string", async () => {
    const { client, last } = fake();
    await client.hosts({ hostsOnly: true });
    expect(last().method).toBe("GET");
    expect(last().query.get("hostsOnly")).toBe("true");
  });

  it("only sends a password when given", async () => {
    const { client, last } = fake();
    await client.unlock("L");
    expect(last().form.has("password")).toBe(false);
    await client.unlock("L", { password: "pw" });
    expect(last().form.get("password")).toBe("pw");
  });

  it("uploads torrents as files[] multipart parts", async () => {
    const { client, last } = fake(ok({ files: [{ magnet: "m", id: 9 }] }));
    await client.uploadTorrentFiles([
      { name: "a.torrent", data: new Uint8Array([100, 101]) },
      { name: "b.torrent", data: new Blob(["x"]) },
    ]);
    expect(last().files).toEqual([
      ["files[]", "a.torrent"],
      ["files[]", "b.torrent"],
    ]);
  });
});

describe("routing", () => {
  const operations = new Set<string>();
  for (const [path, item] of Object.entries(spec.paths)) {
    for (const method of Object.keys(item)) operations.add(`${method.toUpperCase()} ${path}`);
  }

  const calls: [string, (c: AllDebrid) => Promise<unknown>, Record<string, string[]>?][] = [
    ["GET /v4/ping", (c) => c.ping()],
    ["GET /v4.1/pin/get", (c) => c.getPin()],
    ["POST /v4/pin/check", (c) => c.checkPin("P", "C"), { pin: ["P"], check: ["C"] }],
    ["GET /v4/hosts", (c) => c.hosts()],
    ["GET /v4/hosts/domains", (c) => c.hostDomains()],
    ["GET /v4/hosts/priority", (c) => c.hostPriority()],
    ["GET /v4/user", (c) => c.user()],
    ["GET /v4.1/user/hosts", (c) => c.userHosts()],
    ["POST /v4/user/verif", (c) => c.verificationStatus("T"), { token: ["T"] }],
    ["POST /v4/user/verif/resend", (c) => c.resendVerification("T"), { token: ["T"] }],
    ["POST /v4/user/notification/clear", (c) => c.clearNotification("N"), { code: ["N"] }],
    ["POST /v4/link/infos", (c) => c.linkInfos(["a", "b"], { password: "p" }), { "link[]": ["a", "b"], password: ["p"] }],
    ["POST /v4/link/redirector", (c) => c.extractRedirector("L"), { link: ["L"] }],
    ["POST /v4/link/unlock", (c) => c.unlock("L"), { link: ["L"] }],
    ["POST /v4/link/streaming", (c) => c.streamingLink("G", "S"), { id: ["G"], stream: ["S"] }],
    ["POST /v4/link/delayed", (c) => c.delayedLink(42), { id: ["42"] }],
    ["POST /v4/magnet/upload", (c) => c.uploadMagnets(["M"]), { "magnets[]": ["M"] }],
    ["POST /v4/magnet/upload/file", (c) => c.uploadTorrentFiles([{ name: "a", data: new Uint8Array([1]) }])],
    ["POST /v4.1/magnet/status", (c) => c.magnets()],
    ["POST /v4/magnet/files", (c) => c.magnetFiles([1, 2]), { "id[]": ["1", "2"] }],
    ["POST /v4/magnet/delete", (c) => c.deleteMagnet(5), { id: ["5"] }],
    ["POST /v4/magnet/restart", (c) => c.restartMagnets([5, 6]), { "ids[]": ["5", "6"] }],
    ["GET /v4/user/links", (c) => c.savedLinks()],
    ["POST /v4/user/links/save", (c) => c.saveLinks(["L"]), { "links[]": ["L"] }],
    ["POST /v4/user/links/delete", (c) => c.deleteSavedLinks(["L"]), { "links[]": ["L"] }],
    ["GET /v4/user/history", (c) => c.history()],
    ["POST /v4/user/history/delete", (c) => c.purgeHistory()],
    ["GET /v4/voucher/balance", (c) => c.voucherBalance()],
    ["POST /v4/voucher/get", (c) => c.vouchers(30, 2), { duration: ["30"], nb: ["2"] }],
    ["POST /v4/voucher/generate", (c) => c.generateVouchers(90, 1), { duration: ["90"], nb: ["1"] }],
  ];

  it.each(calls)("%s", async (want, call, form) => {
    const answers: Record<string, unknown> = {
      "GET /v4/hosts/priority": { hosts: {} },
      "GET /v4/user": { user: { limitedHostersQuotas: [] } },
      "GET /v4.1/user/hosts": { hosts: {} },
      "POST /v4/link/infos": { infos: [] },
      "POST /v4/link/redirector": { links: [] },
      "POST /v4/magnet/upload": { magnets: [] },
      "POST /v4/magnet/upload/file": { files: [] },
      "POST /v4.1/magnet/status": { magnets: [] },
      "POST /v4/magnet/files": { magnets: [] },
      "GET /v4/user/links": { links: [] },
      "GET /v4/user/history": { links: [] },
      "GET /v4/voucher/balance": { balance: 1 },
    };
    const { client, last } = fake(ok(answers[want] ?? {}));
    await call(client);
    const got = `${last().method} ${last().path}`;
    expect(got).toBe(want);
    expect(operations.has(got)).toBe(true);
    for (const [k, v] of Object.entries(form ?? {})) expect(last().form.getAll(k)).toEqual(v);
  });

  it("covers every operation of the spec", () => {
    expect(new Set(calls.map(([op]) => op))).toEqual(operations);
  });

  it("lists the error codes of the spec", () => {
    expect([...ERROR_CODES].sort()).toEqual([...spec.components.schemas.ErrorCode.enum].sort());
  });
});

describe("responses", () => {
  it("throws API errors with their code", async () => {
    const { client } = fake(apiError("MAGNET_INVALID_ID"));
    const err = await client.deleteMagnet(1).catch((e: unknown) => e);
    expect(err).toBeInstanceOf(AllDebridError);
    expect((err as AllDebridError).is("MAGNET_INVALID_ID")).toBe(true);
    expect(isErrorCode(err, "MAGNET_INVALID_ID")).toBe(true);
    expect(isErrorCode(err, "LINK_DOWN")).toBe(false);
  });

  it("keeps the verification token of AUTH_BLOCKED", async () => {
    const { client } = fake(apiError("AUTH_BLOCKED", { token: "tok" }));
    const err = (await client.user().catch((e: unknown) => e)) as AllDebridError;
    expect(err.code).toBe("AUTH_BLOCKED");
    expect(err.token).toBe("tok");
  });

  it("reports HTTP errors and rate limiting", async () => {
    const { client } = fake(
      { status: 429, body: "Too Many Requests" },
      { status: 404, ...apiError("404") },
    );
    const limited = (await client.ping().catch((e: unknown) => e)) as HttpError;
    expect(limited).toBeInstanceOf(HttpError);
    expect(limited.rateLimited).toBe(true);
    const missing = await client.ping().catch((e: unknown) => e);
    expect(isErrorCode(missing, "404")).toBe(true);
  });

  it("rejects a body that is not JSON", async () => {
    const { client } = fake({ body: "<html>" });
    await expect(client.ping()).rejects.toThrow(/not JSON/);
  });

  it("returns one magnet from the single-object answer", async () => {
    const { client } = fake(ok({ magnets: { id: 7, statusCode: MagnetStatus.Ready, files: [{ n: "a", s: 1, l: "L" }] } }));
    const m = await client.magnet(7);
    expect(m.id).toBe(7);
    expect(m.files?.[0]?.l).toBe("L");
  });

  it("returns a list even when the API sends one object", async () => {
    const { client } = fake(ok({ magnets: { id: 1 } }));
    expect(await client.magnets()).toEqual([{ id: 1 }]);
  });

  it("turns empty quota lists into objects", async () => {
    const { client } = fake(ok({ user: { username: "u", limitedHostersQuotas: [] } }));
    expect((await client.user()).limitedHostersQuotas).toEqual({});
  });
});

describe("helpers", () => {
  it("waits for a PIN", async () => {
    const { client, last } = fake(ok({ activated: false, expires_in: 600 }), ok({ activated: true, expires_in: 590, apikey: "K" }));
    expect(await client.waitForPin({ pin: "ABCD", check: "chk" }, { intervalMs: 1 })).toBe("K");
    expect(last().form.get("pin")).toBe("ABCD");
  });

  it("stops waiting for a PIN when it expires", async () => {
    const { client } = fake(apiError("PIN_EXPIRED"));
    await expect(client.waitForPin({ pin: "A", check: "c" }, { intervalMs: 1 })).rejects.toBeInstanceOf(AllDebridError);
  });

  it("waits for a delayed link", async () => {
    const { client } = fake(ok({ status: 1, time_left: 5 }), ok({ status: 2, time_left: 0, link: "L" }));
    expect(await client.waitDelayedLink("1", { intervalMs: 1 })).toBe("L");
    const failed = fake(ok({ status: 3, time_left: 0 }));
    await expect(failed.client.waitDelayedLink("1", { intervalMs: 1 })).rejects.toThrow(/failed/);
  });

  it("can be aborted", async () => {
    const { client } = fake(ok({ status: 1, time_left: 5 }));
    const controller = new AbortController();
    const waiting = client.waitDelayedLink("1", { intervalMs: 60_000, signal: controller.signal });
    controller.abort(new Error("stop"));
    await expect(waiting).rejects.toThrow("stop");
  });

  it("merges live mode deltas", async () => {
    const { client, requests } = fake(
      ok({ fullsync: true, counter: 1, magnets: [{ id: 1, filename: "a", uploadDate: 10, downloaded: 5 }, { id: 2, filename: "b", uploadDate: 20 }] }),
      ok({ counter: 2, magnets: { id: 1, downloaded: 9 } }),
      ok({ counter: 3, magnets: [{ id: 2, deleted: true }, { id: 3, filename: "c", uploadDate: 30 }] }),
      ok({ counter: 4, magnets: [] }),
      ok({ fullsync: true, counter: 1, magnets: [{ id: 9, uploadDate: 1 }] }),
    );
    const live = client.liveMagnets();
    expect((await live.poll()).map((m) => m.id)).toEqual([2, 1]);
    expect(requests[0]!.form.get("counter")).toBe("0");

    const merged = await live.poll();
    expect(requests[1]!.form.get("counter")).toBe("1");
    expect(requests[1]!.form.get("session")).toBe(requests[0]!.form.get("session"));
    expect(merged.find((m) => m.id === 1)).toMatchObject({ filename: "a", downloaded: 9 });

    expect((await live.poll()).map((m) => m.id)).toEqual([3, 1]);
    expect((await live.poll()).map((m) => m.id)).toEqual([3, 1]);
    expect((await live.poll()).map((m) => m.id)).toEqual([9]);
  });
});
