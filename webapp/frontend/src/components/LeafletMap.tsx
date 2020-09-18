import { Map, Marker, TileLayer } from 'react-leaflet'

import type { FC } from 'react'
import type { Coordinate } from 'types'

interface Props {
  className?: string
  center: Coordinate
  zoom: number
  markerPositions?: Coordinate[]
}

export const LeafletMap: FC<Props> = ({ center, zoom, markerPositions, ...props }) => {
  return (
    <Map
      {...props}
      center={[center.latitude, center.longitude]}
      zoom={zoom}
    >
      <TileLayer
        attribution='&amp;copy <a href=&quot;http://osm.org/copyright&quot;>OpenStreetMap</a> contributors'
        url='https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
      />
      {
        (markerPositions ?? []).map((position, i) => (
          <Marker key={i} position={[position.latitude, position.longitude]} />
        ))
      }
    </Map>
  )
}
