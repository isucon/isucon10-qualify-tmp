<?php

namespace App\Domain;

class Estate
{
    public function getId(): ?int
    {
        return is_null($this->id) ?: (int)$this->id;
    }

    public function getThumbnail(): ?string
    {
        return $this->thumbnail;
    }

    public function getName(): ?string
    {
        return $this->name;
    }

    public function getDescription(): ?string
    {
        return $this->description;
    }

    public function getLatitude(): ?float
    {
        return is_null($this->latitude) ?: (float)$this->latitude;
    }

    public function getLongitude(): ?float
    {
        return is_null($this->longitude) ?: (float)$this->longitude;
    }

    public function getAddress(): ?string
    {
        return $this->address;
    }

    public function getRent(): ?int
    {
        return is_null($this->rent) ?: (int)$this->rent;
    }

    public function getDoorHeight(): ?int
    {
        return is_null($this->door_height) ?: (int)$this->door_height;
    }

    public function getDoorWidth(): ?int
    {
        return is_null($this->door_width) ?: (int)$this->door_width;
    }

    public function getFeatures(): ?string
    {
        return $this->features;
    }

    public function getPopularity(): ?int
    {
        return is_null($this->popularity) ?: (int)$this->popularity;
    }

    public function toArray()
    {
        return [
            'id' => $this->getId(),
            'thumbnail' => $this->getThumbnail(),
            'name' => $this->getName(),
            'description' => $this->getDescription(),
            'latitude' => $this->getLatitude(),
            'longitude' => $this->getLongitude(),
            'address' => $this->getAddress(),
            'rent' => $this->getRent(),
            'doorHeight' => $this->getDoorHeight(),
            'doorWidth' => $this->getDoorWidth(),
            'features' => $this->getFeatures(),
            'popularity' => $this->getPopularity(),
        ];
    }
}