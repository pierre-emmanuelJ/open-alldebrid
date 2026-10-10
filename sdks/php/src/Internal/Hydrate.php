<?php

declare(strict_types=1);

namespace OpenAllDebrid\Internal;

use UnexpectedValueException;

/**
 * Converts decoded JSON values to typed values, for the generated models.
 * Every method throws UnexpectedValueException on a value of the wrong type.
 *
 * @internal
 */
final class Hydrate
{
    /** @param array<array-key, mixed> $data */
    public static function required(array $data, string $key): mixed
    {
        if (!array_key_exists($key, $data)) {
            throw new UnexpectedValueException("missing field \"{$key}\"");
        }
        return $data[$key];
    }

    public static function string(mixed $value): string
    {
        if (!is_string($value)) {
            throw new UnexpectedValueException('expected a string, got ' . get_debug_type($value));
        }
        return $value;
    }

    public static function int(mixed $value): int
    {
        if (is_int($value)) {
            return $value;
        }
        if (is_float($value) && floor($value) === $value) {
            return (int) $value;
        }
        if (is_string($value) && preg_match('/^-?\d+$/', $value) === 1) {
            return (int) $value;
        }
        throw new UnexpectedValueException('expected an integer, got ' . get_debug_type($value));
    }

    public static function float(mixed $value): float
    {
        if (is_int($value) || is_float($value)) {
            return (float) $value;
        }
        throw new UnexpectedValueException('expected a number, got ' . get_debug_type($value));
    }

    public static function bool(mixed $value): bool
    {
        if (!is_bool($value)) {
            throw new UnexpectedValueException('expected a boolean, got ' . get_debug_type($value));
        }
        return $value;
    }

    /** @return array<array-key, mixed> */
    public static function object(mixed $value): array
    {
        if (!is_array($value)) {
            throw new UnexpectedValueException('expected an object, got ' . get_debug_type($value));
        }
        return $value;
    }

    /**
     * @template T
     * @param callable(mixed): T $item
     * @return list<T>
     */
    public static function list(mixed $value, callable $item): array
    {
        if (!is_array($value) || !array_is_list($value)) {
            throw new UnexpectedValueException('expected a list, got ' . get_debug_type($value));
        }
        return array_map($item, $value);
    }

    /**
     * @template T
     * @param callable(mixed): T $item
     * @return array<string, T>
     */
    public static function map(mixed $value, callable $item): array
    {
        $out = [];
        foreach (self::object($value) as $key => $v) {
            $out[(string) $key] = $item($v);
        }
        return $out;
    }

    /** A value the API sends as a string or a number. */
    public static function flexString(mixed $value): string
    {
        if (is_string($value)) {
            return $value;
        }
        if (is_int($value) || is_float($value)) {
            return (string) $value;
        }
        throw new UnexpectedValueException('expected a string or a number, got ' . get_debug_type($value));
    }

    /**
     * One string or a list of strings.
     *
     * @return list<string>
     */
    public static function strings(mixed $value): array
    {
        if (is_string($value)) {
            return [$value];
        }
        return self::list($value, static fn (mixed $v): string => self::string($v));
    }

    /**
     * A list, or a single object to wrap in a list.
     *
     * @return list<mixed>
     */
    public static function oneOrMany(mixed $value): array
    {
        if ($value === null) {
            return [];
        }
        $array = self::object($value);
        return array_is_list($array) ? $array : [$array];
    }

    /**
     * An object, or an empty list standing for an empty object.
     *
     * @return array<array-key, mixed>
     */
    public static function emptyListAsMap(mixed $value): array
    {
        $array = self::object($value);
        if ($array !== [] && array_is_list($array)) {
            throw new UnexpectedValueException('expected an object, got a non-empty list');
        }
        return $array;
    }
}
