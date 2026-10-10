<?php

declare(strict_types=1);

namespace OpenAllDebrid\Tests;

use OpenAllDebrid\Transport;

/** Records requests and answers canned bodies in order. */
final class FakeTransport implements Transport
{
    /** @var list<array{method: string, url: string, headers: array<string, string>, body: ?string}> */
    public array $requests = [];

    /** @param list<array{0: int, 1: string}> $answers */
    public function __construct(private array $answers = [])
    {
    }

    public function answer(int $status, string $body): void
    {
        $this->answers[] = [$status, $body];
    }

    public function ok(string $data): void
    {
        $this->answer(200, '{"status":"success","data":' . $data . '}');
    }

    public function send(string $method, string $url, array $headers, ?string $body): array
    {
        $this->requests[] = ['method' => $method, 'url' => $url, 'headers' => $headers, 'body' => $body];
        return array_shift($this->answers) ?? [200, '{"status":"success","data":{}}'];
    }

    /** @return array{method: string, url: string, headers: array<string, string>, body: ?string} */
    public function last(): array
    {
        $last = end($this->requests);
        if ($last === false) {
            throw new \LogicException('no request');
        }
        return $last;
    }

    public function lastPath(): string
    {
        return (string) parse_url($this->last()['url'], PHP_URL_PATH);
    }

    /**
     * The decoded form body of the last request, keeping repeated keys.
     *
     * @return array<string, list<string>>
     */
    public function lastForm(): array
    {
        $body = $this->last()['body'] ?? '';
        $out = [];
        foreach ($body === '' ? [] : explode('&', $body) as $pair) {
            [$k, $v] = array_pad(explode('=', $pair, 2), 2, '');
            $out[rawurldecode($k)][] = rawurldecode($v);
        }
        return $out;
    }
}
