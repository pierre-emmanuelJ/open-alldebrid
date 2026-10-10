<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class MagnetStatusData
{
    public function __construct(
        /** @var list<Magnet> */
        public readonly array $magnets,
        /** Live mode. The counter to send on the next call. */
        public readonly ?int $counter = null,
        /** Live mode. `true` when the response holds the full list. */
        public readonly ?bool $fullsync = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            magnets: Hydrate::list(Hydrate::oneOrMany(Hydrate::required($data, 'magnets')), static fn (mixed $v): Magnet => Magnet::fromArray(Hydrate::object($v))),
            counter: array_key_exists('counter', $data) && $data['counter'] !== null ? Hydrate::int($data['counter']) : null,
            fullsync: array_key_exists('fullsync', $data) && $data['fullsync'] !== null ? Hydrate::bool($data['fullsync']) : null,
        );
    }
}
