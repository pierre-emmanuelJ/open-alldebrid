<?php

declare(strict_types=1);

namespace OpenAllDebrid\Tests;

use OpenAllDebrid\AllDebrid;
use OpenAllDebrid\AllDebridException;
use OpenAllDebrid\DelayedStatus;
use OpenAllDebrid\ErrorCode;
use PHPUnit\Framework\TestCase;

/** Calls the real API with the public demo keys. Set OFFLINE=1 to skip. */
final class DemoTest extends TestCase
{
    private AllDebrid $client;
    private AllDebrid $anonymous;

    protected function setUp(): void
    {
        if (getenv('OFFLINE') === '1') {
            self::markTestSkipped('OFFLINE=1');
        }
        $this->client = new AllDebrid(AllDebrid::DEMO_PREMIUM_KEY, userAgent: 'open-alldebrid-php-tests');
        $this->anonymous = new AllDebrid(userAgent: 'open-alldebrid-php-tests');
    }

    public function testPublicEndpoints(): void
    {
        $this->anonymous->ping();
        self::assertNotEmpty($this->anonymous->hosts()->hosts);
        self::assertNotEmpty($this->anonymous->hostDomains()->hosts);
        self::assertNotEmpty($this->anonymous->hostPriority());
    }

    public function testPin(): void
    {
        self::assertStringStartsWith('https://', $this->anonymous->getPin()->userUrl);
        self::assertTrue($this->anonymous->checkPin('123B', '664c3ca2635c99f291d28e11ea18e154750bd21a')->activated);
    }

    public function testUser(): void
    {
        self::assertSame('demoUserPremium', $this->client->user()->username);
        self::assertNotEmpty($this->client->userHosts());
        $this->client->clearNotification('NOTIF_CODE');
        self::assertSame('allowed', $this->anonymous->verificationStatus('verificationTokenAllowed3f9eb0db')->verif);
        try {
            $this->anonymous->resendVerification('verificationTokenAlreadySentb0db');
            self::fail('no exception');
        } catch (AllDebridException $e) {
            self::assertSame(ErrorCode::ALREADY_SENT, $e->errorCode);
        }
    }

    public function testLinks(): void
    {
        $infos = $this->client->linkInfos(['https://example.com/file2MB.txt', 'https://example.com/down']);
        self::assertNotEmpty($infos[0]->filename);
        self::assertNotEmpty($this->client->extractRedirector('https://example.com/redir/folder'));
        self::assertNotEmpty($this->client->unlock('https://example.com/streaming')->streaming);
        self::assertNotNull($this->client->unlock('https://example.com/delayed')->delayed);
        self::assertNotEmpty($this->client->streamingLink('demov1bc797', 'hd-720')->link);
        self::assertSame(DelayedStatus::READY, $this->client->delayedLink('1234567')->status);
        try {
            $this->client->unlock('https://example.com/down');
            self::fail('no exception');
        } catch (AllDebridException $e) {
            self::assertSame(ErrorCode::LINK_DOWN, $e->errorCode);
        }
    }

    public function testMagnets(): void
    {
        $uploaded = $this->client->uploadMagnets(['magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707']);
        self::assertNotNull($uploaded[0]->id);
        $this->client->uploadTorrentFiles(['demo.torrent' => 'd4:infod4:name4:demoee']);
        self::assertNotEmpty($this->client->magnets());
        self::assertNotEmpty($this->client->magnet(56789)->files);
        self::assertCount(1, $this->client->magnetFiles([56789]));
        $this->client->deleteMagnet(56789);
        $this->client->restartMagnet(99999);
        self::assertNotEmpty($this->client->liveMagnets()->poll());
    }

    public function testSavedLinks(): void
    {
        self::assertNotEmpty($this->client->savedLinks());
        $this->client->saveLinks(['https://example.com/file2MB.txt']);
        $this->client->deleteSavedLinks(['https://example.com/file2MB.txt']);
        $this->client->history();
        $this->client->purgeHistory();
        $this->addToAssertionCount(1);
    }

    public function testBadKey(): void
    {
        try {
            $this->client->withApiKey('not-a-real-key')->user();
            self::fail('no exception');
        } catch (AllDebridException $e) {
            self::assertSame(ErrorCode::AUTH_BAD_APIKEY, $e->errorCode);
        }
    }
}
