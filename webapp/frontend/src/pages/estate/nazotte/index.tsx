import dynamic from 'next/dynamic'
import { Container, Paper } from '@material-ui/core'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import { Loading } from '../../../components/Loading'

import type { Coordinate } from '@types'
import type { Theme } from '@material-ui/core/styles'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    container: {
      padding: theme.spacing(2),
      width: '100vw',
      height: '100vh',
      flex: 'auto'
    },
    page: {
      width: '100%',
      height: '100%'
    }
  })
)

const NazottePage = () => {
  const classes = useStyles()
  const NazotteMap = dynamic(
    async () => {
      const module = await import('../../../components/NazotteMap')
      return module.NazzoteMap
    },
    { loading: () => <Loading />, ssr: false }
  )
  const estateCoordinate: Coordinate = {
    latitude: 35.67832667,
    longitude: 139.77044378
  }

  return (
    <Container maxWidth={false} className={classes.container}>
      <Paper className={classes.page}>
        <NazotteMap
          center={estateCoordinate}
          zoom={9}
        />
      </Paper>
    </Container>
  )
}

export default NazottePage
