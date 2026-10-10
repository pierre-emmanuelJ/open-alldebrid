<?php

declare(strict_types=1);

namespace OpenAllDebrid\Tests;

use OpenAllDebrid\AllDebrid;
use OpenAllDebrid\AllDebridException;
use OpenAllDebrid\ErrorCode;
use OpenAllDebrid\HttpException;
use OpenAllDebrid\MagnetStatus;
use OpenAllDebrid\Model\Magnet;
use OpenAllDebrid\Model\PinGetData;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;
use RuntimeException;
use UnexpectedValueException;

final class ClientTest extends TestCase
{
    private const SPEC = __DIR__ . '/../../../openapi.yaml';
    private const USER = '{"user":{"username":"u","email":"e","isPremium":true,"isSubscribed":false,"isTrial":false,'
        . '"premiumUntil":0,"lang":"en","preferedDomain":"com","fidelityPoints":0,"limitedHostersQuotas":[],"notifications":[]}}';

    private FakeTransport $transport;
    private AllDebrid $client;

    protected function setUp(): void
    {
        $this->transport = new FakeTransport();
        $this->client = new AllDebrid('test-key', 'https://api.test', 'unit', $this->transport);
    }

    public function testListParametersArePhpStyle(): void
    {
        $this->transport->ok('{"magnets":[{"magnet":"a","id":1},{"magnet":"b","error":{"code":"MAGNET_INVALID_URI","message":"m"}}]}');
        $result = $this->client->uploadMagnets(['a', 'b']);
        $last = $this->transport->last();
        self::assertSame('POST', $last['method']);
        self::assertSame('/v4/magnet/upload', $this->transport->lastPath());
        self::assertSame(['magnets[]' => ['a', 'b']], $this->transport->lastForm());
        self::assertSame('Bearer test-key', $last['headers']['Authorization']);
        self::assertSame('unit', $last['headers']['User-Agent']);
        self::assertSame(1, $result[0]->id);
        self::assertSame('MAGNET_INVALID_URI', $result[1]->error?->code);
    }

    public function testNoKeyNoAuthorization(): void
    {
        $this->transport->ok('{"ping":"pong"}');
        $this->client->withApiKey(null)->ping();
        self::assertArrayNotHasKey('Authorization', $this->transport->last()['headers']);
    }

    public function testGetParametersInQuery(): void
    {
        $this->transport->ok('{"hosts":{}}');
        $this->client->hosts(hostsOnly: true);
        self::assertSame('GET', $this->transport->last()['method']);
        self::assertStringEndsWith('/v4/hosts?hostsOnly=true', $this->transport->last()['url']);
    }

    public function testTorrentFilesAreMultipart(): void
    {
        $this->transport->ok('{"files":[{"magnet":"m","id":9}]}');
        $this->client->uploadTorrentFiles(['a.torrent' => 'd4:infoe', 'b.torrent' => 'x']);
        $last = $this->transport->last();
        self::assertStringStartsWith('multipart/form-data; boundary=', $last['headers']['Content-Type']);
        self::assertSame(2, substr_count((string) $last['body'], 'name="files[]"'));
        self::assertStringContainsString('filename="a.torrent"', (string) $last['body']);
        self::assertStringContainsString('filename="b.torrent"', (string) $last['body']);
    }

    /** @return array<string, array{0: string, 1: callable(AllDebrid): mixed, 2: array<string, list<string>>, 3: string}> */
    public static function routes(): array
    {
        $routes = [
            ['GET /v4/ping', static fn (AllDebrid $c) => $c->ping(), [], '{}'],
            ['GET /v4.1/pin/get', static fn (AllDebrid $c) => $c->getPin(), [], '{"pin":"P","check":"C","expires_in":600,"user_url":"u","base_url":"b"}'],
            ['POST /v4/pin/check', static fn (AllDebrid $c) => $c->checkPin('P', 'C'), ['pin' => ['P'], 'check' => ['C']], '{"activated":false,"expires_in":600}'],
            ['GET /v4/hosts', static fn (AllDebrid $c) => $c->hosts(), [], '{"hosts":{}}'],
            ['GET /v4/hosts/domains', static fn (AllDebrid $c) => $c->hostDomains(), [], '{"hosts":[]}'],
            ['GET /v4/hosts/priority', static fn (AllDebrid $c) => $c->hostPriority(), [], '{"hosts":{}}'],
            ['GET /v4/user', static fn (AllDebrid $c) => $c->user(), [], self::USER],
            ['GET /v4.1/user/hosts', static fn (AllDebrid $c) => $c->userHosts(), [], '{"hosts":{}}'],
            ['POST /v4/user/verif', static fn (AllDebrid $c) => $c->verificationStatus('T'), ['token' => ['T']], '{"verif":"waiting"}'],
            ['POST /v4/user/verif/resend', static fn (AllDebrid $c) => $c->resendVerification('T'), ['token' => ['T']], '{}'],
            ['POST /v4/user/notification/clear', static fn (AllDebrid $c) => $c->clearNotification('N'), ['code' => ['N']], '{}'],
            ['POST /v4/link/infos', static fn (AllDebrid $c) => $c->linkInfos(['a', 'b'], 'p'), ['link[]' => ['a', 'b'], 'password' => ['p']], '{"infos":[]}'],
            ['POST /v4/link/redirector', static fn (AllDebrid $c) => $c->extractRedirector('L'), ['link' => ['L']], '{"links":[]}'],
            ['POST /v4/link/unlock', static fn (AllDebrid $c) => $c->unlock('L'), ['link' => ['L']], '{}'],
            ['POST /v4/link/streaming', static fn (AllDebrid $c) => $c->streamingLink('G', 'S'), ['id' => ['G'], 'stream' => ['S']], '{}'],
            ['POST /v4/link/delayed', static fn (AllDebrid $c) => $c->delayedLink('42'), ['id' => ['42']], '{"status":1,"time_left":5}'],
            ['POST /v4/magnet/upload', static fn (AllDebrid $c) => $c->uploadMagnets(['M']), ['magnets[]' => ['M']], '{"magnets":[]}'],
            ['POST /v4/magnet/upload/file', static fn (AllDebrid $c) => $c->uploadTorrentFiles(['a' => 'x']), [], '{"files":[]}'],
            ['POST /v4.1/magnet/status', static fn (AllDebrid $c) => $c->magnets(), [], '{"magnets":[]}'],
            ['POST /v4/magnet/files', static fn (AllDebrid $c) => $c->magnetFiles([1, 2]), ['id[]' => ['1', '2']], '{"magnets":[]}'],
            ['POST /v4/magnet/delete', static fn (AllDebrid $c) => $c->deleteMagnet(5), ['id' => ['5']], '{}'],
            ['POST /v4/magnet/restart', static fn (AllDebrid $c) => $c->restartMagnets([5, 6]), ['ids[]' => ['5', '6']], '{}'],
            ['GET /v4/user/links', static fn (AllDebrid $c) => $c->savedLinks(), [], '{"links":[]}'],
            ['POST /v4/user/links/save', static fn (AllDebrid $c) => $c->saveLinks(['L']), ['links[]' => ['L']], '{}'],
            ['POST /v4/user/links/delete', static fn (AllDebrid $c) => $c->deleteSavedLinks(['L']), ['links[]' => ['L']], '{}'],
            ['GET /v4/user/history', static fn (AllDebrid $c) => $c->history(), [], '{"links":[]}'],
            ['POST /v4/user/history/delete', static fn (AllDebrid $c) => $c->purgeHistory(), [], '{}'],
            ['GET /v4/voucher/balance', static fn (AllDebrid $c) => $c->voucherBalance(), [], '{"balance":1.5}'],
            ['POST /v4/voucher/get', static fn (AllDebrid $c) => $c->vouchers(30, 2), ['duration' => ['30'], 'nb' => ['2']], '{"codes":[]}'],
            ['POST /v4/voucher/generate', static fn (AllDebrid $c) => $c->generateVouchers(90, 1), ['duration' => ['90'], 'nb' => ['1']], '{"codes":[]}'],
        ];
        $out = [];
        foreach ($routes as $route) {
            $out[$route[0]] = $route;
        }
        return $out;
    }

    /**
     * @param callable(AllDebrid): mixed $call
     * @param array<string, list<string>> $form
     */
    #[DataProvider('routes')]
    public function testRouting(string $want, callable $call, array $form, string $data): void
    {
        $this->transport->ok($data);
        $call($this->client);
        $got = $this->transport->last()['method'] . ' ' . $this->transport->lastPath();
        self::assertSame($want, $got);
        self::assertContains($got, self::specOperations());
        if ($form !== []) {
            self::assertSame($form, $this->transport->lastForm());
        }
    }

    public function testEveryOperationIsRouted(): void
    {
        $routed = array_keys(self::routes());
        sort($routed);
        $ops = self::specOperations();
        sort($ops);
        self::assertSame($ops, $routed);
    }

    public function testErrorCodesMatchTheSpec(): void
    {
        $spec = (string) file_get_contents(self::SPEC);
        $block = explode("\n    ErrorCode:\n", $spec, 2)[1];
        $block = preg_split('/\n    \S/', $block, 2)[0] ?? '';
        preg_match_all('/^        - "?([A-Z0-9_]+)"?$/m', $block, $m);
        $codes = $m[1];
        $known = ErrorCode::ALL;
        sort($codes);
        sort($known);
        self::assertSame($codes, $known);
    }

    /** @return list<string> "METHOD /path" of every operation of the spec. */
    private static function specOperations(): array
    {
        $ops = [];
        $path = null;
        foreach (explode("\n", (string) file_get_contents(self::SPEC)) as $line) {
            if ($line === 'components:') {
                break;
            }
            if (preg_match('#^  (/\S+):$#', $line, $m) === 1) {
                $path = $m[1];
            } elseif ($path !== null && preg_match('/^    (get|post):$/', $line, $m) === 1) {
                $ops[] = strtoupper($m[1]) . ' ' . $path;
            }
        }
        return $ops;
    }

    public function testApiErrorsCarryTheirCode(): void
    {
        $this->transport->answer(200, '{"status":"error","error":{"code":"AUTH_BLOCKED","message":"m","token":"tok"}}');
        try {
            $this->client->user();
            self::fail('no exception');
        } catch (AllDebridException $e) {
            self::assertSame(ErrorCode::AUTH_BLOCKED, $e->errorCode);
            self::assertSame('tok', $e->token);
            self::assertSame('m', $e->apiMessage);
        }
    }

    public function testHttpErrors(): void
    {
        $this->transport->answer(429, 'Too Many Requests');
        $this->transport->answer(404, '{"status":"error","error":{"code":"404","message":"Endpoint doesn\'t exist"}}');
        $this->transport->answer(200, '<html>');
        try {
            $this->client->ping();
            self::fail('no exception');
        } catch (HttpException $e) {
            self::assertTrue($e->isRateLimited());
        }
        try {
            $this->client->ping();
            self::fail('no exception');
        } catch (HttpException $e) {
            self::assertSame(ErrorCode::NOT_FOUND, $e->apiError?->errorCode);
        }
        $this->expectException(UnexpectedValueException::class);
        $this->client->ping();
    }

    public function testFlexibleShapes(): void
    {
        $this->transport->ok('{"magnets":{"id":7,"statusCode":4,"files":[{"n":"a","s":1,"l":"L"}]}}');
        $magnet = $this->client->magnet(7);
        self::assertSame(7, $magnet->id);
        self::assertTrue(MagnetStatus::isReady($magnet));
        self::assertSame('L', $magnet->files[0]->l ?? null);

        $this->transport->ok(self::USER);
        self::assertSame([], $this->client->user()->limitedHostersQuotas);

        $this->transport->ok('{"link":"","delayed":123,"streaming":[{"id":"mp3","quality":"mp3"},{"id":"hd","quality":1080}]}');
        $unlocked = $this->client->unlock('L');
        self::assertSame('123', $unlocked->delayed);
        self::assertSame(['mp3', '1080'], array_map(static fn ($s) => $s->quality, $unlocked->streaming ?? []));

        $this->transport->ok('{"hosts":{"a":{"name":"a","type":"free","domains":[],"regexp":"x"}}}');
        self::assertSame(['x'], $this->client->hosts()->hosts['a']->regexp);
    }

    public function testWrongTypesAreRejected(): void
    {
        $this->transport->ok('{"user":{"username":1}}');
        $this->expectException(UnexpectedValueException::class);
        $this->client->user();
    }

    public function testWaitForPin(): void
    {
        $this->transport->ok('{"activated":false,"expires_in":600}');
        $this->transport->ok('{"activated":true,"expires_in":590,"apikey":"K"}');
        $pin = new PinGetData('ABCD', 'chk', 600, 'u', 'b');
        self::assertSame('K', $this->client->waitForPin($pin, 0.001));
        self::assertSame(['pin' => ['ABCD'], 'check' => ['chk']], $this->transport->lastForm());
    }

    public function testWaitDelayedLink(): void
    {
        $this->transport->ok('{"status":1,"time_left":5}');
        $this->transport->ok('{"status":2,"time_left":0,"link":"L"}');
        self::assertSame('L', $this->client->waitDelayedLink('1', 0.001));
        $this->transport->ok('{"status":3,"time_left":0}');
        $this->expectException(RuntimeException::class);
        $this->client->waitDelayedLink('1', 0.001);
    }

    /**
     * @param list<Magnet> $magnets
     * @return list<int>
     */
    private static function ids(array $magnets): array
    {
        return array_map(static fn (Magnet $m): int => $m->id, $magnets);
    }

    public function testLiveMagnets(): void
    {
        $this->transport->ok('{"fullsync":true,"counter":1,"magnets":[{"id":1,"filename":"a","uploadDate":10,"downloaded":5},{"id":2,"filename":"b","uploadDate":20}]}');
        $this->transport->ok('{"counter":2,"magnets":{"id":1,"downloaded":9}}');
        $this->transport->ok('{"counter":3,"magnets":[{"id":2,"deleted":true},{"id":3,"filename":"c","uploadDate":30}]}');
        $this->transport->ok('{"counter":4,"magnets":[]}');
        $this->transport->ok('{"fullsync":true,"counter":1,"magnets":[{"id":9,"uploadDate":1}]}');
        $live = $this->client->liveMagnets();

        self::assertSame([2, 1], self::ids($live->poll()));
        $first = $this->transport->lastForm();
        self::assertSame(['0'], $first['counter']);

        $merged = $live->poll();
        $second = $this->transport->lastForm();
        self::assertSame(['1'], $second['counter']);
        self::assertSame($first['session'], $second['session']);
        self::assertSame('a', $merged[1]->filename);
        self::assertSame(9, $merged[1]->downloaded);

        $afterDelete = self::ids($live->poll());
        self::assertSame([3, 1], $afterDelete);
        $unchanged = self::ids($live->poll());
        self::assertSame([3, 1], $unchanged, 'an empty delta changed the state');
        $resynced = self::ids($live->poll());
        self::assertSame([9], $resynced);
    }
}
