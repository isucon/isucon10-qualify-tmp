import { useEffect, useState, useCallback, useRef } from 'react'
import { useRouter } from 'next/router'
import dynamic from 'next/dynamic'
import {
  Paper,
  Container,
  Box,
  TextField,
  Button
} from '@material-ui/core'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import { Loading } from '../../../components/Loading'

import type { FC } from 'react'
import type { Estate, Coordinate } from '@types'
import type { Theme } from '@material-ui/core/styles'
import ErrorPage from 'next/error'

const usePageStyles = makeStyles((theme: Theme) =>
  createStyles({
    page: {
      margin: theme.spacing(2),
      padding: theme.spacing(4)
    }
  })
)

const useEstateDetailStyles = makeStyles((theme: Theme) =>
  createStyles({
    column: {
      marginTop: theme.spacing(4),
      marginBottom: theme.spacing(4)
    },
    row: {
      '&>*': {
        margin: theme.spacing(2)
      }
    },
    thumbnailContainer: {
      height: 270
    },
    thumbnail: {
      height: '100%'
    },
    map: {
      width: '100%',
      height: 360
    }
  })
)

interface Props {
  estate: Estate
}

const EstateDetail: FC<Props> = ({ estate }) => {
  const classes = useEstateDetailStyles()
  const LeafletMap = dynamic(
    async () => {
      const module = await import('../../../components/LeafletMap')
      return module.LeafletMap
    },
    { ssr: false }
  )
  const estateCoordinate: Coordinate = {
    latitude: estate.latitude,
    longitude: estate.longitude
  }

  const emailInputRef = useRef<HTMLInputElement>(null)
  const [submitResult, setSubmitResult] = useState<string>('')

  const onSubmit = useCallback(async () => {
    const EMAIL_REGEXP = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/
    if (!EMAIL_REGEXP.test(emailInputRef.current?.value ?? '')) {
      setSubmitResult('Invalid email address format.')
      return
    }

    await fetch(`/api/estate/req_doc/${estate.id}`, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      },
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify({ email: emailInputRef.current?.value })
    })
      .then(async response => response.status.toString() + (await response.text()))
      .then(setSubmitResult)
      .catch(error => setSubmitResult(error.message))
  }, [estate.id])

  return (
    <>
      <Box width={1} className={`${classes.column} ${classes.thumbnailContainer}`} display='flex' justifyContent='center'>
        <img src={estate.thumbnail} alt='物件画像' className={classes.thumbnail} />
      </Box>

      {estate.id && (
        <Box width={1} className={classes.row} display='flex' alignItems='center'>
          <h2 style={{ wordBreak: 'keep-all' }}>資料請求:</h2>

          <TextField label='Email address' inputRef={emailInputRef} fullWidth />
          <Button variant='contained' color='primary' onClick={onSubmit}> Submit </Button>
          {submitResult && <p> {submitResult} </p>}
        </Box>
      )}

      <Box width={1} className={classes.column}>
        <h2>この物件について</h2>

        <p>名前: {estate.name}</p>
        <p>説明: {estate.description}</p>
        <p>賃料: {estate.rent}円</p>
        <p>住所: {estate.address}</p>
        <LeafletMap
          className={classes.map}
          center={estateCoordinate}
          zoom={10}
          markerPositions={[estateCoordinate]}
        />
      </Box>

      <Box width={1} className={classes.column}>
        <h3>ドアのサイズ</h3>
        <ul>
          <li>縦 (cm) : {estate.doorHeight}</li>
          <li>横 (cm) : {estate.doorWidth}</li>
        </ul>
      </Box>

      <Box width={1} className={classes.column}>
        <h3>こだわり条件:</h3>
        {
          estate.features === '' ? 'なし' : (
            estate.features.split(',').map((feature, i) => (
              <p key={i}>{feature}</p>
            ))
          )
        }
      </Box>
    </>
  )
}

const EstateDetailPage = () => {
  const [estate, setEstate] = useState<Estate | null>(null)
  const [statusCode, setStatusCode] = useState(200)
  const router = useRouter()
  const id = Array.isArray(router.query.id) ? router.query.id[0] : router.query.id

  const classes = usePageStyles()

  useEffect(() => {
    if (!id) return
    fetch(`/api/estate/${id.toString()}`, { mode: 'cors' })
      .then(async response => {
        if (response.status !== 200) setStatusCode(response.status)
        return await response.json()
      })
      .then(estate => setEstate(estate as Estate))
      .catch(console.error)
  }, [id])

  if (!id) return <ErrorPage statusCode={404} title='Page /estate/detail is required id query like /estate/detail?id=1' />

  if (statusCode !== 200) return <ErrorPage statusCode={statusCode} />

  return (
    <Paper className={classes.page}>
      <Container maxWidth='md'>
        {estate ? (
          <EstateDetail estate={estate} />
        ) : (
          <Loading />
        )}
      </Container>
    </Paper>
  )
}

export default EstateDetailPage
