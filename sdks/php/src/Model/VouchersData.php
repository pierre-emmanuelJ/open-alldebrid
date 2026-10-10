<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class VouchersData
{
    public function __construct(
        /** @var list<string> */
        public readonly array $codes,
        /** `true` when fewer vouchers than requested were available. */
        public readonly ?bool $partialList = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            codes: Hydrate::list(Hydrate::required($data, 'codes'), static fn (mixed $v): string => Hydrate::string($v)),
            partialList: array_key_exists('partialList', $data) && $data['partialList'] !== null ? Hydrate::bool($data['partialList']) : null,
        );
    }
}
