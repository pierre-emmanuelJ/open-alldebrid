<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class HostsData
{
    public function __construct(
        /** @var array<string, Service> */
        public readonly array $hosts,
        /** @var array<string, Service>|null */
        public readonly ?array $streams = null,
        /** @var array<string, Service>|null */
        public readonly ?array $redirectors = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            hosts: Hydrate::map(Hydrate::required($data, 'hosts'), static fn (mixed $v): Service => Service::fromArray(Hydrate::object($v))),
            streams: array_key_exists('streams', $data) && $data['streams'] !== null ? Hydrate::map($data['streams'], static fn (mixed $v): Service => Service::fromArray(Hydrate::object($v))) : null,
            redirectors: array_key_exists('redirectors', $data) && $data['redirectors'] !== null ? Hydrate::map($data['redirectors'], static fn (mixed $v): Service => Service::fromArray(Hydrate::object($v))) : null,
        );
    }
}
