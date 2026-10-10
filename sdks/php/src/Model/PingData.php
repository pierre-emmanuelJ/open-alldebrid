<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class PingData
{
    public function __construct(
        public readonly string $ping,
        /** API major version. */
        public readonly ?string $version = null,
        /** Only on `/v4.1/ping`. */
        public readonly ?bool $subversion = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            ping: Hydrate::string(Hydrate::required($data, 'ping')),
            version: array_key_exists('version', $data) && $data['version'] !== null ? Hydrate::string($data['version']) : null,
            subversion: array_key_exists('subversion', $data) && $data['subversion'] !== null ? Hydrate::bool($data['subversion']) : null,
        );
    }
}
