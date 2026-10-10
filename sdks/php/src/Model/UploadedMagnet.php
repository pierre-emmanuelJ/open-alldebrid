<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

/** An uploaded magnet or torrent, or a per-item `error`. */
final class UploadedMagnet
{
    public function __construct(
        /** The magnet as sent (for torrent files, the magnet built from the file). */
        public readonly ?string $magnet = null,
        /** Torrent file uploads with an `error` only. The file name sent. */
        public readonly ?string $file = null,
        /** The info hash. */
        public readonly ?string $hash = null,
        /** The torrent name, or `noname` when unknown. */
        public readonly ?string $name = null,
        public readonly ?string $filenameOriginal = null,
        /** Total size in bytes. */
        public readonly ?int $size = null,
        /** Whether the files are already available. */
        public readonly ?bool $ready = null,
        /** The magnet ID, for `/v4.1/magnet/status`. */
        public readonly ?int $id = null,
        public readonly ?ApiError $error = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            magnet: array_key_exists('magnet', $data) && $data['magnet'] !== null ? Hydrate::string($data['magnet']) : null,
            file: array_key_exists('file', $data) && $data['file'] !== null ? Hydrate::string($data['file']) : null,
            hash: array_key_exists('hash', $data) && $data['hash'] !== null ? Hydrate::string($data['hash']) : null,
            name: array_key_exists('name', $data) && $data['name'] !== null ? Hydrate::string($data['name']) : null,
            filenameOriginal: array_key_exists('filename_original', $data) && $data['filename_original'] !== null ? Hydrate::string($data['filename_original']) : null,
            size: array_key_exists('size', $data) && $data['size'] !== null ? Hydrate::int($data['size']) : null,
            ready: array_key_exists('ready', $data) && $data['ready'] !== null ? Hydrate::bool($data['ready']) : null,
            id: array_key_exists('id', $data) && $data['id'] !== null ? Hydrate::int($data['id']) : null,
            error: array_key_exists('error', $data) && $data['error'] !== null ? ApiError::fromArray(Hydrate::object($data['error'])) : null,
        );
    }
}
