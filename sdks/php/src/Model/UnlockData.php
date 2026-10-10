<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class UnlockData
{
    public function __construct(
        /** The direct download link. Empty when `delayed` is set. */
        public readonly ?string $link = null,
        public readonly ?string $host = null,
        public readonly ?string $filename = null,
        /** File size in bytes. */
        public readonly ?int $filesize = null,
        /** The generation ID, needed by `/v4/link/streaming`. */
        public readonly ?string $id = null,
        public readonly ?string $hostDomain = null,
        /** Unused. */
        public readonly ?bool $paws = null,
        /**
         * Available qualities for streaming links; empty otherwise.
         * @var list<Stream>|null
         */
        public readonly ?array $streaming = null,
        /** Set when the link is still being generated; poll `/v4/link/delayed` with it. */
        public readonly ?string $delayed = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            link: array_key_exists('link', $data) && $data['link'] !== null ? Hydrate::string($data['link']) : null,
            host: array_key_exists('host', $data) && $data['host'] !== null ? Hydrate::string($data['host']) : null,
            filename: array_key_exists('filename', $data) && $data['filename'] !== null ? Hydrate::string($data['filename']) : null,
            filesize: array_key_exists('filesize', $data) && $data['filesize'] !== null ? Hydrate::int($data['filesize']) : null,
            id: array_key_exists('id', $data) && $data['id'] !== null ? Hydrate::string($data['id']) : null,
            hostDomain: array_key_exists('hostDomain', $data) && $data['hostDomain'] !== null ? Hydrate::string($data['hostDomain']) : null,
            paws: array_key_exists('paws', $data) && $data['paws'] !== null ? Hydrate::bool($data['paws']) : null,
            streaming: array_key_exists('streaming', $data) && $data['streaming'] !== null ? Hydrate::list($data['streaming'], static fn (mixed $v): Stream => Stream::fromArray(Hydrate::object($v))) : null,
            delayed: array_key_exists('delayed', $data) && $data['delayed'] !== null ? Hydrate::flexString($data['delayed']) : null,
        );
    }
}
