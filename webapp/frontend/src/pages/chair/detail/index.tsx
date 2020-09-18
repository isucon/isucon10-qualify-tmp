import { useEffect, useState, useCallback, useRef } from 'react'
import { useRouter } from 'next/router'
import {
  Paper,
  Container,
  Box,
  TextField,
  Button
} from '@material-ui/core'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import { Loading } from '../../../components/Loading'
import { EstateCard } from '../../../components/EstateCard'
import ErrorPage from 'next/error'

import type { FC } from 'react'
import type { Chair, Estate } from '@types'
import type { Theme } from '@material-ui/core/styles'

const usePageStyles = makeStyles((theme: Theme) =>
  createStyles({
    page: {
      margin: theme.spacing(2),
      padding: theme.spacing(4)
    }
  })
)

const useChairDetailStyles = makeStyles((theme: Theme) =>
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
    cards: {
      display: 'flex',
      flexWrap: 'wrap',
      justifyContent: 'space-between'
    }
  })
)

interface Props {
  chair: Chair
  lowPricedEstates: Estate[]
}

const ChairDetail: FC<Props> = ({ chair, lowPricedEstates }) => {
  const classes = useChairDetailStyles()

  const emailInputRef = useRef<HTMLInputElement>(null)
  const [submitResult, setSubmitResult] = useState<string>('')

  const onSubmit = useCallback(async () => {
    const EMAIL_REGEXP = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/
    if (!EMAIL_REGEXP.test(emailInputRef.current?.value ?? '')) {
      setSubmitResult('Invalid email address format.')
      return
    }

    await fetch(`/api/chair/buy/${chair.id}`, {
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
  }, [chair.id])

  return (
    <>
      <Box width={1} className={`${classes.column} ${classes.thumbnailContainer}`} display='flex' justifyContent='center'>
        <img src={chair.thumbnail} alt='イスの画像' className={classes.thumbnail} />
      </Box>

      {chair.id && (
        <Box width={1} className={classes.row} display='flex' alignItems='center'>
          <h2 style={{ wordBreak: 'keep-all' }}>購入:</h2>

          <TextField label='Email address' inputRef={emailInputRef} fullWidth />
          <Button variant='contained' color='primary' onClick={onSubmit}> Submit </Button>
          {submitResult && <p> {submitResult} </p>}
        </Box>
      )}

      <Box width={1} className={classes.column}>
        <h2>このイスについて</h2>

        <p>名前: {chair.name}</p>
        <p>説明: {chair.description}</p>
        <p>値段: {chair.price}円</p>
        <p>色: {chair.color}</p>
        <p>種類: {chair.kind}</p>
      </Box>

      <Box width={1} className={classes.column}>
        <h3>イスのサイズ</h3>
        <ul>
          <li>縦 (cm) : {chair.height}</li>
          <li>横 (cm) : {chair.width}</li>
          <li>奥 (cm) : {chair.depth}</li>
        </ul>
      </Box>

      <Box width={1} className={classes.column}>
        <h3>こだわり条件:</h3>
        {
          chair.features === '' ? 'なし' : (
            chair.features.split(',').map((feature, i) => (
              <p key={i}>{feature}</p>
            ))
          )
        }
      </Box>

      <Box width={1} className={classes.column}>
        <h3>このイスにオススメの物件:</h3>
        <Box width={1} className={classes.cards}>
          {lowPricedEstates.map((estate, i) => <EstateCard key={i} estate={estate} />)}
        </Box>
      </Box>
    </>
  )
}

const ChairDetailPage = () => {
  const [chair, setChair] = useState<Chair | null>(null)
  const [statusCode, setStatusCode] = useState(200)
  const [lowPricedEstates, setLowPricedEstates] = useState<Estate[] | null>(null)
  const router = useRouter()
  const id = Array.isArray(router.query.id) ? router.query.id[0] : router.query.id

  const classes = usePageStyles()

  useEffect(() => {
    if (!id) return

    fetch(`/api/chair/${id.toString()}`, { mode: 'cors' })
      .then(async response => {
        if (response.status !== 200) setStatusCode(response.status)
        return await response.json()
      })
      .then(chair => setChair(chair as Chair))
      .catch(error => { throw error })

    fetch(`/api/recommended_estate/${id.toString()}`, { mode: 'cors' })
      .then(async response => await response.json())
      .then(json => setLowPricedEstates(json.estates as Estate[]))
      .catch(error => { throw error })
  }, [id])

  if (!id) return <ErrorPage statusCode={404} title='Page /chair/detail is required id query like /chair/detail?id=1' />

  if (statusCode !== 200) return <ErrorPage statusCode={statusCode} />

  return (
    <Paper className={classes.page}>
      <Container maxWidth='md'>
        {chair && lowPricedEstates ? (
          <ChairDetail chair={chair} lowPricedEstates={lowPricedEstates} />
        ) : (
          <Loading />
        )}
      </Container>
    </Paper>
  )
}

export default ChairDetailPage
