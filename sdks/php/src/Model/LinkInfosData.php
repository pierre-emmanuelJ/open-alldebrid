<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class LinkInfosData
{
    public function __construct(
        /** @var list<LinkInfo> */
        public readonly array $infos,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            infos: Hydrate::list(Hydrate::required($data, 'infos'), static fn (mixed $v): LinkInfo => LinkInfo::fromArray(Hydrate::object($v))),
        );
    }
}
