<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

/**
 * A magnet of the account. In a live mode delta, only `id` and the
 * fields that changed are present.
 */
final class Magnet
{
    public function __construct(
        public readonly int $id,
        public readonly ?string $filename = null,
        /** Total size in bytes. */
        public readonly ?int $size = null,
        /** The info hash. */
        public readonly ?string $hash = null,
        /** The status in plain English. */
        public readonly ?string $status = null,
        /**
         * `0` queued, `1` downloading, `2` compressing or moving, `3` uploading,
         * `4` ready. `5` to `15` are errors: `5` upload failed, `6` unpacking error,
         * `7` no peer, `8` too big, `9` internal error, `10` took too long,
         * `11` files removed from the host, `12` and `13` processing failed,
         * `14` tracker error, `15` no peer.
         */
        public readonly ?int $statusCode = null,
        /** Bytes downloaded so far (while processing). */
        public readonly ?int $downloaded = null,
        /** Bytes uploaded so far (while processing). */
        public readonly ?int $uploaded = null,
        /** Seeders (while downloading). */
        public readonly ?int $seeders = null,
        /** Bytes per second (while downloading). */
        public readonly ?int $downloadSpeed = null,
        /** Bytes per second (while uploading). */
        public readonly ?int $uploadSpeed = null,
        /** Progress percentage while processing. */
        public readonly ?float $processingPerc = null,
        /** Unix time the magnet was added. */
        public readonly ?int $uploadDate = null,
        /** Unix time the magnet became ready, or 0. */
        public readonly ?int $completionDate = null,
        /** `m` when added as a magnet, `b` when added as a torrent file. */
        public readonly ?string $type = null,
        public readonly ?bool $notified = null,
        public readonly ?int $version = null,
        /** Number of files, once ready. */
        public readonly ?int $nbLinks = null,
        /**
         * Only when a single magnet is requested with `id`.
         * @var list<FileNode>|null
         */
        public readonly ?array $files = null,
        /** Live mode deltas only. The magnet was deleted; drop it from the local state. */
        public readonly ?bool $deleted = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            id: Hydrate::int(Hydrate::required($data, 'id')),
            filename: array_key_exists('filename', $data) && $data['filename'] !== null ? Hydrate::string($data['filename']) : null,
            size: array_key_exists('size', $data) && $data['size'] !== null ? Hydrate::int($data['size']) : null,
            hash: array_key_exists('hash', $data) && $data['hash'] !== null ? Hydrate::string($data['hash']) : null,
            status: array_key_exists('status', $data) && $data['status'] !== null ? Hydrate::string($data['status']) : null,
            statusCode: array_key_exists('statusCode', $data) && $data['statusCode'] !== null ? Hydrate::int($data['statusCode']) : null,
            downloaded: array_key_exists('downloaded', $data) && $data['downloaded'] !== null ? Hydrate::int($data['downloaded']) : null,
            uploaded: array_key_exists('uploaded', $data) && $data['uploaded'] !== null ? Hydrate::int($data['uploaded']) : null,
            seeders: array_key_exists('seeders', $data) && $data['seeders'] !== null ? Hydrate::int($data['seeders']) : null,
            downloadSpeed: array_key_exists('downloadSpeed', $data) && $data['downloadSpeed'] !== null ? Hydrate::int($data['downloadSpeed']) : null,
            uploadSpeed: array_key_exists('uploadSpeed', $data) && $data['uploadSpeed'] !== null ? Hydrate::int($data['uploadSpeed']) : null,
            processingPerc: array_key_exists('processingPerc', $data) && $data['processingPerc'] !== null ? Hydrate::float($data['processingPerc']) : null,
            uploadDate: array_key_exists('uploadDate', $data) && $data['uploadDate'] !== null ? Hydrate::int($data['uploadDate']) : null,
            completionDate: array_key_exists('completionDate', $data) && $data['completionDate'] !== null ? Hydrate::int($data['completionDate']) : null,
            type: array_key_exists('type', $data) && $data['type'] !== null ? Hydrate::string($data['type']) : null,
            notified: array_key_exists('notified', $data) && $data['notified'] !== null ? Hydrate::bool($data['notified']) : null,
            version: array_key_exists('version', $data) && $data['version'] !== null ? Hydrate::int($data['version']) : null,
            nbLinks: array_key_exists('nbLinks', $data) && $data['nbLinks'] !== null ? Hydrate::int($data['nbLinks']) : null,
            files: array_key_exists('files', $data) && $data['files'] !== null ? Hydrate::list($data['files'], static fn (mixed $v): FileNode => FileNode::fromArray(Hydrate::object($v))) : null,
            deleted: array_key_exists('deleted', $data) && $data['deleted'] !== null ? Hydrate::bool($data['deleted']) : null,
        );
    }
}
