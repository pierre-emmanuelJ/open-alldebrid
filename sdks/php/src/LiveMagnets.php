<?php

declare(strict_types=1);

namespace OpenAllDebrid;

use OpenAllDebrid\Internal\Hydrate;
use OpenAllDebrid\Model\Magnet;

/** Live mode state: the merged list of magnets, updated by poll(). */
final class LiveMagnets
{
    public readonly int $session;
    private int $counter = 0;
    /** @var array<int, array<array-key, mixed>> */
    private array $state = [];

    public function __construct(private readonly AllDebrid $client)
    {
        $this->session = random_int(0, 2 ** 31 - 1);
    }

    /**
     * Applies the changes since the last call and returns every magnet, newest first.
     *
     * @return list<Magnet>
     * @phpstan-impure
     */
    public function poll(): array
    {
        $data = Hydrate::object($this->client->call('POST', '/v4.1/magnet/status', [
            ['session', (string) $this->session],
            ['counter', (string) $this->counter],
        ]));
        if (($data['fullsync'] ?? false) === true) {
            $this->state = [];
        }
        foreach (Hydrate::oneOrMany($data['magnets'] ?? null) as $change) {
            $change = Hydrate::object($change);
            $id = Hydrate::int(Hydrate::required($change, 'id'));
            if (($change['deleted'] ?? false) === true) {
                unset($this->state[$id]);
            } else {
                $this->state[$id] = array_merge($this->state[$id] ?? [], $change);
            }
        }
        $this->counter = Hydrate::int($data['counter'] ?? 0);

        $magnets = array_map(static fn (array $m): Magnet => Magnet::fromArray($m), array_values($this->state));
        usort($magnets, static fn (Magnet $a, Magnet $b): int => [$b->uploadDate ?? 0, $b->id] <=> [$a->uploadDate ?? 0, $a->id]);
        return $magnets;
    }
}
