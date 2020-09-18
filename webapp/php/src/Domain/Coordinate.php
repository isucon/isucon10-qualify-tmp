<?php

namespace App\Domain;

class Coordinate
{
    public float $latitude;
    public float $longitude;

    public function __construct(float $latitude, float $longitude)
    {
        $this->latitude = $latitude;
        $this->longitude = $longitude;
    }

    public static function createFromJson(array $coordinate): Coordinate
    {
        return new Coordinate($coordinate['latitude'], $coordinate['longitude']);
    }

    /**
     * @param Coordinate[]
     */
    public static function toText(array $coordinates): string
    {
        return sprintf("'POLYGON((%s))'", implode(',', array_map(
            function(Coordinate $coordinate) {
                return sprintf('%f %f', $coordinate->latitude, $coordinate->longitude);
            },
            $coordinates
        )));
    }
}