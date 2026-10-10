<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class GeneratedVouchersData
{
    public function __construct(
        /** @var list<string> */
        public readonly array $codes,
        public readonly ?float $pricePerVoucher = null,
        /** Total price of the request. */
        public readonly ?float $total = null,
        /** The balance left after the generation. */
        public readonly ?float $balance = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            codes: Hydrate::list(Hydrate::required($data, 'codes'), static fn (mixed $v): string => Hydrate::string($v)),
            pricePerVoucher: array_key_exists('pricePerVoucher', $data) && $data['pricePerVoucher'] !== null ? Hydrate::float($data['pricePerVoucher']) : null,
            total: array_key_exists('total', $data) && $data['total'] !== null ? Hydrate::float($data['total']) : null,
            balance: array_key_exists('balance', $data) && $data['balance'] !== null ? Hydrate::float($data['balance']) : null,
        );
    }
}
