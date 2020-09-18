<?php

namespace App\Domain;

class BoundingBox
{
    public Coordinate $topLeftCorner;
    public Coordinate $bottomRightCorner;

    public function __construct(Coordinate $topLeftCorner, Coordinate $bottomRightCorner)
    {
        $this->topLeftCorner = $topLeftCorner;
        $this->bottomRightCorner = $bottomRightCorner;
    }

    /**
     * @param Coordinate[] $coordinates
     */
    public static function createFromCordinates(array $coordinates): BoundingBox
    {
        $boundingBox = new BoundingBox(
            new Coordinate($coordinates[0]->latitude, $coordinates[0]->longitude),
            new Coordinate($coordinates[0]->latitude, $coordinates[0]->longitude),
        );

        foreach ($coordinates as $coordinate) {
            if ($boundingBox->topLeftCorner->latitude > $coordinate->latitude) {
                $boundingBox->topLeftCorner->latitude = $coordinate->latitude;
            }
            if ($boundingBox->topLeftCorner->longitude > $coordinate->longitude) {
                $boundingBox->topLeftCorner->longitude = $coordinate->longitude;
            }
            if ($boundingBox->bottomRightCorner->latitude < $coordinate->latitude) {
                $boundingBox->bottomRightCorner->latitude = $coordinate->latitude;
            }
            if ($boundingBox->bottomRightCorner->longitude < $coordinate->longitude) {
                $boundingBox->bottomRightCorner->longitude = $coordinate->longitude;
            }
        }

        return $boundingBox;
    }
}