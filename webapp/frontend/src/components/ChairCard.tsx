import Link from 'next/link'
import {
  Card,
  CardActionArea,
  CardMedia,
  CardContent
} from '@material-ui/core'
import { makeStyles, createStyles } from '@material-ui/core/styles'

import type { FC } from 'react'
import type { Chair } from '@types'

interface Props {
  chair: Chair
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

export const ChairCard: FC<Props> = ({ chair }) => {
  const classes = useStyles()
  return (
    <Link href={`/chair/detail?id=${chair.id}`}>
      <Card className={classes.cardContainer}>
        <CardActionArea component='div' disableRipple className={classes.card}>
          <CardMedia
            className={classes.cardMedia}
            image={chair.thumbnail}
            title={chair.name}
          />
          <CardContent>
            <h3> {chair.name} </h3>
            <p> 価格: {chair.price}円 </p>
            <p> 種類: {chair.kind} </p>
          </CardContent>
        </CardActionArea>
      </Card>
    </Link>
  )
}
