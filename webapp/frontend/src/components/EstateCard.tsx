import Link from 'next/link'
import {
  Card,
  CardActionArea,
  CardMedia,
  CardContent
} from '@material-ui/core'
import { makeStyles, createStyles } from '@material-ui/core/styles'

import type { FC } from 'react'
import type { Estate } from '@types'

interface Props {
  estate: Estate
}

const useStyles = makeStyles(theme =>
  createStyles({
    cardContainer: {
      margin: theme.spacing(1),
      width: 270,
      minWidth: 270,
      height: 270
    },
    card: {
      width: '100%',
      height: '100%'
    },
    cardMedia: {
      width: '100%',
      height: 120
    }
  })
)

export const EstateCard: FC<Props> = ({ estate }) => {
  const classes = useStyles()
  return (
    <Link href={`/estate/detail?id=${estate.id}`}>
      <Card className={classes.cardContainer}>
        <CardActionArea component='div' disableRipple className={classes.card}>
          <CardMedia
            className={classes.cardMedia}
            image={estate.thumbnail}
            title={estate.name}
          />
          <CardContent>
            <h3>{estate.name}</h3>
            <p>住所: {estate.address}</p>
            <p>家賃: {estate.rent}円</p>
          </CardContent>
        </CardActionArea>
      </Card>
    </Link>
  )
}
