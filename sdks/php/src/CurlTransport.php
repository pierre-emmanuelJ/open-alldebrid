<?php

declare(strict_types=1);

namespace OpenAllDebrid;

use RuntimeException;

/** The default transport, on ext-curl. */
final class CurlTransport implements Transport
{
    public function __construct(private readonly int $timeoutSeconds = 30)
    {
    }

    public function send(string $method, string $url, array $headers, ?string $body): array
    {
        if ($method === '') {
            throw new RuntimeException('empty HTTP method');
        }
        $handle = curl_init($url);
        if ($handle === false) {
            throw new RuntimeException('curl_init failed');
        }
        $lines = [];
        foreach ($headers as $name => $value) {
            $lines[] = "{$name}: {$value}";
        }
        curl_setopt_array($handle, [
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_HTTPHEADER => $lines,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT => $this->timeoutSeconds,
            CURLOPT_FOLLOWLOCATION => false,
        ]);
        if ($body !== null) {
            curl_setopt($handle, CURLOPT_POSTFIELDS, $body);
        }
        $response = curl_exec($handle);
        if (!is_string($response)) {
            $error = curl_error($handle);
            throw new RuntimeException("AllDebrid {$method} request failed: {$error}");
        }
        return [curl_getinfo($handle, CURLINFO_RESPONSE_CODE), $response];
    }
}
