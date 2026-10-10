<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class StreamingData
{
    public function __construct(
        /** The download link; absent when `delayed` is set. */
        public readonly ?string $link = null,
        public readonly ?string $filename = null,
        public readonly ?int $filesize = null,
        /** Set when the link is still being generated; poll `/v4/link/delayed` with it. */
        public readonly ?string $delayed = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            link: array_key_exists('link', $data) && $data['link'] !== null ? Hydrate::string($data['link']) : null,
            filename: array_key_exists('filename', $data) && $data['filename'] !== null ? Hydrate::string($data['filename']) : null,
            filesize: array_key_exists('filesize', $data) && $data['filesize'] !== null ? Hydrate::int($data['filesize']) : null,
            delayed: array_key_exists('delayed', $data) && $data['delayed'] !== null ? Hydrate::flexString($data['delayed']) : null,
        );
    }
}
