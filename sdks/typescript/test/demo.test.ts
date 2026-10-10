import { describe, expect, it } from "vitest";
import { AllDebrid, DEMO_KEYS, DelayedStatus, isErrorCode } from "../src/index.js";

// Calls the real API with the public demo keys. Set OFFLINE=1 to skip.
describe.skipIf(process.env["OFFLINE"] === "1")("demo API", () => {
  const client = new AllDebrid({ apiKey: DEMO_KEYS.premium, userAgent: "open-alldebrid-ts-tests" });
  const anonymous = new AllDebrid({ userAgent: "open-alldebrid-ts-tests" });

  it("public endpoints", async () => {
    await anonymous.ping();
    const hosts = await anonymous.hosts();
    expect(Object.keys(hosts.hosts).length).toBeGreaterThan(0);
    expect((await anonymous.hostDomains()).hosts.length).toBeGreaterThan(0);
    expect(Object.keys(await anonymous.hostPriority()).length).toBeGreaterThan(0);
  });

  it("PIN flow", async () => {
    const pin = await anonymous.getPin();
    expect(pin.user_url).toMatch(/^https:\/\//);
    const state = await anonymous.checkPin("123B", "664c3ca2635c99f291d28e11ea18e154750bd21a");
    expect(state.activated).toBe(true);
  });

  it("user", async () => {
    const user = await client.user();
    expect(user.username).toBe("demoUserPremium");
    expect(Object.keys(await client.userHosts()).length).toBeGreaterThan(0);
    await client.clearNotification("NOTIF_CODE");
    expect((await anonymous.verificationStatus("verificationTokenAllowed3f9eb0db")).verif).toBe("allowed");
    const resent = await anonymous.resendVerification("verificationTokenAlreadySentb0db").catch((e: unknown) => e);
    expect(isErrorCode(resent, "ALREADY_SENT")).toBe(true);
  });

  it("links", async () => {
    const infos = await client.linkInfos(["https://example.com/file2MB.txt", "https://example.com/down"]);
    expect(infos[0]?.filename).toBeTruthy();
    expect(infos.length).toBeGreaterThan(0);
    expect((await client.extractRedirector("https://example.com/redir/folder")).length).toBeGreaterThan(0);
    const streaming = await client.unlock("https://example.com/streaming");
    expect(streaming.streaming?.length).toBeGreaterThan(0);
    const delayed = await client.unlock("https://example.com/delayed");
    expect(delayed.delayed).toBeTruthy();
    expect(isErrorCode(await client.unlock("https://example.com/down").catch((e: unknown) => e), "LINK_DOWN")).toBe(true);
    expect((await client.streamingLink("demov1bc797", "hd-720")).link).toBeTruthy();
    expect((await client.delayedLink("1234567")).status).toBe(DelayedStatus.Ready);
  });

  it("magnets", async () => {
    const uploaded = await client.uploadMagnets(["magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707"]);
    expect(uploaded[0]?.id).toBeTruthy();
    await client.uploadTorrentFiles([{ name: "demo.torrent", data: new TextEncoder().encode("d4:infod4:name4:demoee") }]);
    expect((await client.magnets()).length).toBeGreaterThan(0);
    const one = await client.magnet(56789);
    expect(one.files?.length).toBeGreaterThan(0);
    expect((await client.magnetFiles([56789])).length).toBe(1);
    await client.deleteMagnet(56789);
    await client.restartMagnet(99999);
    expect((await client.liveMagnets().poll()).length).toBeGreaterThan(0);
  });

  it("saved links and history", async () => {
    expect((await client.savedLinks()).length).toBeGreaterThan(0);
    await client.saveLinks(["https://example.com/file2MB.txt"]);
    await client.deleteSavedLinks(["https://example.com/file2MB.txt"]);
    await client.history();
    await client.purgeHistory();
  });

  it("rejects a bad key", async () => {
    const err = await client.withApiKey("not-a-real-key").user().catch((e: unknown) => e);
    expect(isErrorCode(err, "AUTH_BAD_APIKEY")).toBe(true);
  });
});
