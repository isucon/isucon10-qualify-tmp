<?php

namespace App\Domain;

class Chair
{
    public function getId(): ?int
    {
        return is_null($this->id) ?: (int)$this->id;
    }

    public function getName(): ?string
    {
        return $this->name;
    }

    public function getDescription(): ?string
    {
        return $this->description;
    }

    public function getThumbnail(): ?string
    {
        return $this->thumbnail;
    }

    public function getPrice(): ?int
    {
        return is_null($this->price) ?: (int)$this->price;
    }

    public function getHeight(): ?int
    {
        return is_null($this->height) ?: (int)$this->height;
    }

    public function getWidth(): ?int
    {
        return is_null($this->width) ?: (int)$this->width;
    }

    public function getDepth(): ?int
    {
        return is_null($this->depth) ?: (int)$this->depth;
    }

    public function getColor(): ?string
    {
        return $this->color;
    }

    public function getFeatures(): ?string
    {
        return $this->features;
    }

    public function getKind(): ?string
    {
        return $this->kind;
    }

    public function getPopularity(): ?int
    {
        return is_null($this->popularity) ?: (int)$this->popularity;
    }

    public function getStock(): ?int
    {
        return is_null($this->stock) ?: (int)$this->stock;
    }

    public function toArray()
    {
        return [
            'id' => $this->getId(),
            'name' => $this->getName(),
            'description' => $this->getDescription(),
            'thumbnail' => $this->getThumbnail(),
            'price' => $this->getPrice(),
            'height' => $this->getHeight(),
            'width' => $this->getWidth(),
            'depth' => $this->getDepth(),
            'color' => $this->getColor(),
            'features' => $this->getFeatures(),
            'kind' => $this->getKind(),
            'popularity' => $this->getPopularity(),
            'stock' => $this->getStock(),
        ];
    }
}