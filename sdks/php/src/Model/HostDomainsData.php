<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class HostDomainsData
{
    public function __construct(
        /** @var list<string> */
        public readonly array $hosts,
        /** @var list<string>|null */
        public readonly ?array $streams = null,
        /** @var list<string>|null */
        public readonly ?array $redirectors = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            hosts: Hydrate::list(Hydrate::required($data, 'hosts'), static fn (mixed $v): string => Hydrate::string($v)),
            streams: array_key_exists('streams', $data) && $data['streams'] !== null ? Hydrate::list($data['streams'], static fn (mixed $v): string => Hydrate::string($v)) : null,
            redirectors: array_key_exists('redirectors', $data) && $data['redirectors'] !== null ? Hydrate::list($data['redirectors'], static fn (mixed $v): string => Hydrate::string($v)) : null,
        );
    }
}
