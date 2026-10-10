<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class LinkInfo
{
    public function __construct(
        /** The requested link. */
        public readonly string $link,
        public readonly ?string $filename = null,
        /** File size in bytes. */
        public readonly ?int $size = null,
        public readonly ?string $host = null,
        /** The host's main domain. */
        public readonly ?string $hostDomain = null,
        public readonly ?ApiError $error = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            link: Hydrate::string(Hydrate::required($data, 'link')),
            filename: array_key_exists('filename', $data) && $data['filename'] !== null ? Hydrate::string($data['filename']) : null,
            size: array_key_exists('size', $data) && $data['size'] !== null ? Hydrate::int($data['size']) : null,
            host: array_key_exists('host', $data) && $data['host'] !== null ? Hydrate::string($data['host']) : null,
            hostDomain: array_key_exists('hostDomain', $data) && $data['hostDomain'] !== null ? Hydrate::string($data['hostDomain']) : null,
            error: array_key_exists('error', $data) && $data['error'] !== null ? ApiError::fromArray(Hydrate::object($data['error'])) : null,
        );
    }
}
