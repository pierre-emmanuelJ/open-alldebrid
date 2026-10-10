<?php

declare(strict_types=1);

namespace OpenAllDebrid;

use OpenAllDebrid\Internal\Hydrate;
use OpenAllDebrid\Model\DelayedData;
use OpenAllDebrid\Model\GeneratedVouchersData;
use OpenAllDebrid\Model\HostDomainsData;
use OpenAllDebrid\Model\HostPriorityData;
use OpenAllDebrid\Model\HostsData;
use OpenAllDebrid\Model\LinkInfo;
use OpenAllDebrid\Model\LinkInfosData;
use OpenAllDebrid\Model\Magnet;
use OpenAllDebrid\Model\MagnetFiles;
use OpenAllDebrid\Model\MagnetFilesData;
use OpenAllDebrid\Model\MagnetStatusData;
use OpenAllDebrid\Model\PinCheckData;
use OpenAllDebrid\Model\PinGetData;
use OpenAllDebrid\Model\RedirectorData;
use OpenAllDebrid\Model\RestartedMagnet;
use OpenAllDebrid\Model\RestartMagnetsData;
use OpenAllDebrid\Model\SavedLink;
use OpenAllDebrid\Model\SavedLinksData;
use OpenAllDebrid\Model\Service;
use OpenAllDebrid\Model\StreamingData;
use OpenAllDebrid\Model\UnlockData;
use OpenAllDebrid\Model\UploadedMagnet;
use OpenAllDebrid\Model\UploadMagnetsData;
use OpenAllDebrid\Model\UploadTorrentFilesData;
use OpenAllDebrid\Model\User;
use OpenAllDebrid\Model\UserData;
use OpenAllDebrid\Model\UserHostsData;
use OpenAllDebrid\Model\VerificationData;
use OpenAllDebrid\Model\VoucherBalanceData;
use OpenAllDebrid\Model\VouchersData;
use RuntimeException;
use UnexpectedValueException;

/**
 * A client for the AllDebrid API.
 *
 *     $client = new AllDebrid($apiKey, userAgent: 'my-app');
 *     $user = $client->user();
 *
 * API errors throw AllDebridException, other HTTP statuses HttpException.
 */
final class AllDebrid
{
    public const DEFAULT_BASE_URL = 'https://api.alldebrid.com';
    public const DEFAULT_USER_AGENT = 'open-alldebrid-php';

    /** Public demo keys: the API answers them with static sample data. */
    public const DEMO_PREMIUM_KEY = 'staticDemoApikeyPrem';
    public const DEMO_TRIAL_KEY = 'staticDemoApikeyTria';
    public const DEMO_FREE_KEY = 'staticDemoApikeyFree';

    private readonly string $baseUrl;
    private readonly Transport $transport;

    public function __construct(
        private readonly ?string $apiKey = null,
        string $baseUrl = self::DEFAULT_BASE_URL,
        private readonly string $userAgent = self::DEFAULT_USER_AGENT,
        ?Transport $transport = null,
    ) {
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->transport = $transport ?? new CurlTransport();
    }

    /** A client with another API key, sharing this one's transport. */
    public function withApiKey(?string $apiKey): self
    {
        return new self($apiKey, $this->baseUrl, $this->userAgent, $this->transport);
    }

    /**
     * Calls the API and returns the "data" of a successful answer.
     *
     * @param list<array{0: string, 1: string}> $params
     * @internal
     */
    public function call(string $method, string $path, array $params = [], ?Multipart $multipart = null): mixed
    {
        $headers = ['Accept' => 'application/json', 'User-Agent' => $this->userAgent];
        if ($this->apiKey !== null && $this->apiKey !== '') {
            $headers['Authorization'] = "Bearer {$this->apiKey}";
        }
        $url = $this->baseUrl . $path;
        $body = null;
        if ($multipart !== null) {
            $headers['Content-Type'] = $multipart->contentType();
            $body = $multipart->body();
        } elseif ($method === 'GET') {
            if ($params !== []) {
                $url .= '?' . self::encode($params);
            }
        } else {
            // List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
            $headers['Content-Type'] = 'application/x-www-form-urlencoded';
            $body = self::encode($params);
        }

        [$status, $text] = $this->transport->send($method, $url, $headers, $body);
        $envelope = json_decode($text, true);
        $apiError = is_array($envelope) ? self::apiError($envelope['error'] ?? null) : null;
        if ($status !== 200) {
            throw new HttpException($status, $apiError);
        }
        if (!is_array($envelope)) {
            throw new UnexpectedValueException("AllDebrid {$method} {$path}: the response is not a JSON object");
        }
        if (($envelope['status'] ?? null) !== 'success') {
            throw $apiError ?? new UnexpectedValueException("AllDebrid {$method} {$path}: error without details");
        }
        return $envelope['data'] ?? null;
    }

    /** @param list<array{0: string, 1: string}> $params */
    private static function encode(array $params): string
    {
        return implode('&', array_map(
            static fn (array $kv): string => rawurlencode($kv[0]) . '=' . rawurlencode($kv[1]),
            $params,
        ));
    }

    private static function apiError(mixed $error): ?AllDebridException
    {
        if (!is_array($error)) {
            return null;
        }
        $token = $error['token'] ?? null;
        return new AllDebridException(
            is_scalar($error['code'] ?? null) ? (string) $error['code'] : '',
            is_scalar($error['message'] ?? null) ? (string) $error['message'] : '',
            is_scalar($token) ? (string) $token : null,
        );
    }

    /**
     * @param list<array{0: string, 1: string}> $params
     * @return array<array-key, mixed>
     */
    private function data(string $method, string $path, array $params = [], ?Multipart $multipart = null): array
    {
        return Hydrate::object($this->call($method, $path, $params, $multipart));
    }

    /**
     * @param iterable<string> $values
     * @return list<array{0: string, 1: string}>
     */
    private static function repeat(string $name, iterable $values): array
    {
        $out = [];
        foreach ($values as $value) {
            $out[] = [$name, (string) $value];
        }
        return $out;
    }

    // -------------------------------------------------------------- general

    /** Checks that the API is up. */
    public function ping(): void
    {
        $this->call('GET', '/v4/ping');
    }

    // ------------------------------------------------------------------ pin

    /** Starts a PIN authentication. Show $pin->pin or $pin->userUrl, then waitForPin(). */
    public function getPin(): PinGetData
    {
        return PinGetData::fromArray($this->data('GET', '/v4.1/pin/get'));
    }

    /** The state of a PIN authentication; apikey is set once activated. */
    public function checkPin(string $pin, string $check): PinCheckData
    {
        return PinCheckData::fromArray($this->data('POST', '/v4/pin/check', [['pin', $pin], ['check', $check]]));
    }

    /** Polls until the user confirms the PIN and returns the new API key. */
    public function waitForPin(PinGetData $pin, float $intervalSeconds = 5.0): string
    {
        while (true) {
            $state = $this->checkPin($pin->pin, $pin->check);
            if ($state->activated && $state->apikey !== null && $state->apikey !== '') {
                return $state->apikey;
            }
            usleep((int) ($intervalSeconds * 1_000_000));
        }
    }

    // ---------------------------------------------------------------- hosts

    /** Supported hosts, streaming sites and redirectors. Cache the result. */
    public function hosts(bool $hostsOnly = false): HostsData
    {
        return HostsData::fromArray($this->data('GET', '/v4/hosts', $hostsOnly ? [['hostsOnly', 'true']] : []));
    }

    /** Every supported domain. */
    public function hostDomains(): HostDomainsData
    {
        return HostDomainsData::fromArray($this->data('GET', '/v4/hosts/domains'));
    }

    /**
     * Host names ranked from the most open (1) to the most restricted.
     *
     * @return array<string, int>
     */
    public function hostPriority(): array
    {
        return HostPriorityData::fromArray($this->data('GET', '/v4/hosts/priority'))->hosts;
    }

    // ----------------------------------------------------------------- user

    /** The account information. */
    public function user(): User
    {
        return UserData::fromArray($this->data('GET', '/v4/user'))->user;
    }

    /**
     * The hosts this user can use now, with their quotas.
     *
     * @return array<string, Service>
     */
    public function userHosts(bool $hostsOnly = false): array
    {
        $params = $hostsOnly ? [['hostsOnly', 'true']] : [];
        return UserHostsData::fromArray($this->data('GET', '/v4.1/user/hosts', $params))->hosts;
    }

    /** Follows a new-location login verification (token from an AUTH_BLOCKED error). */
    public function verificationStatus(string $token): VerificationData
    {
        return VerificationData::fromArray($this->data('POST', '/v4/user/verif', [['token', $token]]));
    }

    /** Sends the verification email again (once per token). */
    public function resendVerification(string $token): void
    {
        $this->call('POST', '/v4/user/verif/resend', [['token', $token]]);
    }

    /** Clears one of the codes in User::$notifications. */
    public function clearNotification(string $code): void
    {
        $this->call('POST', '/v4/user/notification/clear', [['code', $code]]);
    }

    // ---------------------------------------------------------------- links

    /**
     * File name and size of each link, in order; unreadable links have an error.
     *
     * @param list<string> $links
     * @return list<LinkInfo>
     */
    public function linkInfos(array $links, ?string $password = null): array
    {
        $params = self::repeat('link[]', $links);
        if ($password !== null) {
            $params[] = ['password', $password];
        }
        return LinkInfosData::fromArray($this->data('POST', '/v4/link/infos', $params))->infos;
    }

    /**
     * Links behind a redirector; they only work with linkInfos() and unlock().
     *
     * @return list<string>
     */
    public function extractRedirector(string $link): array
    {
        return RedirectorData::fromArray($this->data('POST', '/v4/link/redirector', [['link', $link]]))->links;
    }

    /**
     * Unlocks a hoster link. When delayed is set, wait with waitDelayedLink();
     * when streaming lists qualities, choose one with streamingLink().
     */
    public function unlock(string $link, ?string $password = null): UnlockData
    {
        $params = [['link', $link]];
        if ($password !== null) {
            $params[] = ['password', $password];
        }
        return UnlockData::fromArray($this->data('POST', '/v4/link/unlock', $params));
    }

    /** The link of one quality: $id from unlock(), $stream from one of its streaming entries. */
    public function streamingLink(string $id, string $stream): StreamingData
    {
        return StreamingData::fromArray($this->data('POST', '/v4/link/streaming', [['id', $id], ['stream', $stream]]));
    }

    /** The status of a delayed link. */
    public function delayedLink(string $id): DelayedData
    {
        return DelayedData::fromArray($this->data('POST', '/v4/link/delayed', [['id', $id]]));
    }

    /** Polls a delayed link until it is ready and returns the download link. */
    public function waitDelayedLink(string $id, float $intervalSeconds = 5.0): string
    {
        while (true) {
            $state = $this->delayedLink($id);
            if ($state->status === DelayedStatus::READY && $state->link !== null && $state->link !== '') {
                return $state->link;
            }
            if ($state->status === DelayedStatus::FAILED) {
                throw new RuntimeException("AllDebrid: delayed link {$id} failed");
            }
            usleep((int) ($intervalSeconds * 1_000_000));
        }
    }

    // -------------------------------------------------------------- magnets

    /**
     * Adds magnet URIs or info hashes. Results are in order; invalid ones have an error.
     *
     * @param list<string> $magnets
     * @return list<UploadedMagnet>
     */
    public function uploadMagnets(array $magnets): array
    {
        $params = self::repeat('magnets[]', $magnets);
        return UploadMagnetsData::fromArray($this->data('POST', '/v4/magnet/upload', $params))->magnets;
    }

    /**
     * Adds .torrent files, given as file name => content. Invalid files have an error.
     *
     * @param array<string, string> $files
     * @return list<UploadedMagnet>
     */
    public function uploadTorrentFiles(array $files): array
    {
        $multipart = new Multipart();
        foreach ($files as $name => $content) {
            $multipart->addFile('files[]', (string) $name, $content, 'application/x-bittorrent');
        }
        $data = $this->data('POST', '/v4/magnet/upload/file', [], $multipart);
        return UploadTorrentFilesData::fromArray($data)->files;
    }

    /**
     * Every magnet of the account.
     *
     * @return list<Magnet>
     */
    public function magnets(): array
    {
        return MagnetStatusData::fromArray($this->data('POST', '/v4.1/magnet/status'))->magnets;
    }

    /** One magnet, with its files. */
    public function magnet(int $id): Magnet
    {
        $found = MagnetStatusData::fromArray($this->data('POST', '/v4.1/magnet/status', [['id', (string) $id]]))->magnets;
        if (count($found) !== 1) {
            throw new UnexpectedValueException("AllDebrid: magnet {$id}: got " . count($found) . ' magnets');
        }
        return $found[0];
    }

    /**
     * The file tree of each magnet, in order; unknown IDs have an error.
     *
     * @param list<int> $ids
     * @return list<MagnetFiles>
     */
    public function magnetFiles(array $ids): array
    {
        $params = self::repeat('id[]', array_map('strval', $ids));
        return MagnetFilesData::fromArray($this->data('POST', '/v4/magnet/files', $params))->magnets;
    }

    /** Deletes a magnet. */
    public function deleteMagnet(int $id): void
    {
        $this->call('POST', '/v4/magnet/delete', [['id', (string) $id]]);
    }

    /** Restarts a failed magnet. */
    public function restartMagnet(int $id): void
    {
        $this->call('POST', '/v4/magnet/restart', [['id', (string) $id]]);
    }

    /**
     * Restarts failed magnets; each result has its own error.
     *
     * @param list<int> $ids
     * @return list<RestartedMagnet>
     */
    public function restartMagnets(array $ids): array
    {
        $params = self::repeat('ids[]', array_map('strval', $ids));
        return RestartMagnetsData::fromArray($this->data('POST', '/v4/magnet/restart', $params))->magnets ?? [];
    }

    /** Follows the magnets with live mode: each poll() only transfers changes. */
    public function liveMagnets(): LiveMagnets
    {
        return new LiveMagnets($this);
    }

    // ---------------------------------------------------------- saved links

    /**
     * The links the user saved.
     *
     * @return list<SavedLink>
     */
    public function savedLinks(): array
    {
        return SavedLinksData::fromArray($this->data('GET', '/v4/user/links'))->links;
    }

    /** @param list<string> $links Links to save for later. */
    public function saveLinks(array $links): void
    {
        $this->call('POST', '/v4/user/links/save', self::repeat('links[]', $links));
    }

    /** @param list<string> $links Saved links to delete. */
    public function deleteSavedLinks(array $links): void
    {
        $this->call('POST', '/v4/user/links/delete', self::repeat('links[]', $links));
    }

    /**
     * Links unlocked in the last 3 days, when history logging is enabled.
     *
     * @return list<SavedLink>
     */
    public function history(): array
    {
        return SavedLinksData::fromArray($this->data('GET', '/v4/user/history'))->links;
    }

    /** Deletes the whole recent links history. */
    public function purgeHistory(): void
    {
        $this->call('POST', '/v4/user/history/delete');
    }

    // ------------------------------------------------------------- reseller

    /** The reseller balance. Reseller accounts only. */
    public function voucherBalance(): float
    {
        return VoucherBalanceData::fromArray($this->data('GET', '/v4/voucher/balance'))->balance;
    }

    /** Takes vouchers (15, 30, 90, 180 or 365 days) from the reseller stock. */
    public function vouchers(int $durationDays, int $count): VouchersData
    {
        $params = [['duration', (string) $durationDays], ['nb', (string) $count]];
        return VouchersData::fromArray($this->data('POST', '/v4/voucher/get', $params));
    }

    /** Generates new vouchers from the reseller balance. */
    public function generateVouchers(int $durationDays, int $count): GeneratedVouchersData
    {
        $params = [['duration', (string) $durationDays], ['nb', (string) $count]];
        return GeneratedVouchersData::fromArray($this->data('POST', '/v4/voucher/generate', $params));
    }
}
