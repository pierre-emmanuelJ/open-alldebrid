<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class SavedLink
{
    public function __construct(
        public readonly string $link,
        public readonly ?string $filename = null,
        public readonly ?int $size = null,
        /** Unix time the link was saved or unlocked. */
        public readonly ?int $date = null,
        public readonly ?string $host = null,
        /** History only. The download link that was generated. */
        public readonly ?string $linkDl = null,
        /** History only. The error the unlock gave, if any. */
        public readonly ?string $error = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            link: Hydrate::string(Hydrate::required($data, 'link')),
            filename: array_key_exists('filename', $data) && $data['filename'] !== null ? Hydrate::string($data['filename']) : null,
            size: array_key_exists('size', $data) && $data['size'] !== null ? Hydrate::int($data['size']) : null,
            date: array_key_exists('date', $data) && $data['date'] !== null ? Hydrate::int($data['date']) : null,
            host: array_key_exists('host', $data) && $data['host'] !== null ? Hydrate::string($data['host']) : null,
            linkDl: array_key_exists('link_dl', $data) && $data['link_dl'] !== null ? Hydrate::string($data['link_dl']) : null,
            error: array_key_exists('error', $data) && $data['error'] !== null ? Hydrate::string($data['error']) : null,
        );
    }
}
