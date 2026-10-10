<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

/** One quality of a streaming link. */
final class Stream
{
    public function __construct(
        /** The `stream` value for `/v4/link/streaming`. */
        public readonly string $id,
        public readonly ?string $ext = null,
        /** Height in pixels (e.g. `1080`), or a label such as `mp3`. */
        public readonly ?string $quality = null,
        public readonly ?int $filesize = null,
        public readonly ?string $proto = null,
        public readonly ?string $name = null,
        /** Audio-only streams. Total bitrate. */
        public readonly ?float $tb = null,
        /** Audio-only streams. Audio bitrate. */
        public readonly ?int $abr = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            id: Hydrate::string(Hydrate::required($data, 'id')),
            ext: array_key_exists('ext', $data) && $data['ext'] !== null ? Hydrate::string($data['ext']) : null,
            quality: array_key_exists('quality', $data) && $data['quality'] !== null ? Hydrate::flexString($data['quality']) : null,
            filesize: array_key_exists('filesize', $data) && $data['filesize'] !== null ? Hydrate::int($data['filesize']) : null,
            proto: array_key_exists('proto', $data) && $data['proto'] !== null ? Hydrate::string($data['proto']) : null,
            name: array_key_exists('name', $data) && $data['name'] !== null ? Hydrate::string($data['name']) : null,
            tb: array_key_exists('tb', $data) && $data['tb'] !== null ? Hydrate::float($data['tb']) : null,
            abr: array_key_exists('abr', $data) && $data['abr'] !== null ? Hydrate::int($data['abr']) : null,
        );
    }
}
