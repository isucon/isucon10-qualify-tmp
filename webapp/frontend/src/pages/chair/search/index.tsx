import { useEffect, useState } from 'react'
import Link from 'next/link'
import {
  Paper,
  Container,
  Box,
  Button,
  CircularProgress,
  Card,
  CardContent,
  CardMedia,
  CardActionArea
} from '@material-ui/core'
import { Pagination } from '@material-ui/lab'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import { Loading } from '../../../components/Loading'
import { RangeForm } from '../../../components/RangeForm'
import { RadioButtonForm } from '../../../components/RadioButtonForm'
import { CheckboxForm } from '../../../components/CheckboxForm'

import type { FC } from 'react'
import type { ChairSearchCondition, ChairSearchParams, ChairSearchResponse } from '@types'

const ESTATE_COUNTS_PER_PAGE = 20

interface ChairSearchProps {
  chairSearchCondition: ChairSearchCondition
}

const useChairSearchStyles = makeStyles(theme =>
  createStyles({
    page: {
      margin: theme.spacing(2),
      padding: theme.spacing(4)
    },
    search: {
      display: 'flex',
      flexDirection: 'column',
      marginTop: theme.spacing(4),
      marginBottom: theme.spacing(4),
      '&>*': {
        margin: theme.spacing(1)
      }
    },
    row: {
      '&>*': {
        margin: theme.spacing(2)
      }
    },
    card: {
      width: '100%',
      height: 270,
      marginTop: theme.spacing(2),
      marginBottom: theme.spacing(2)
    },
    cardActionArea: {
      height: 270,
      display: 'flex',
      alignItems: 'flex-start',
      justifyContent: 'flex-start'
    },
    cardMedia: {
      width: 360,
      height: 270
    },
    cardContent: {
      width: 'fit-content',
      marginLeft: theme.spacing(1)
    }
  })
)

const searchChair = async (params: ChairSearchParams) => {
  const urlSearchParams = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    urlSearchParams.append(key, value.toString())
  }
  const response = await fetch(`/api/chair/search?${urlSearchParams.toString()}`, { mode: 'cors' })
  const json = await response.json()
  return json as ChairSearchResponse
}

const ChairSearch: FC<ChairSearchProps> = ({ chairSearchCondition }) => {
  const classes = useChairSearchStyles()

  const [priceRangeId, setPriceRangeId] = useState('')
  const [heightRangeId, setHeightRangeId] = useState('')
  const [widthRangeId, setWidthRangeId] = useState('')
  const [depthRangeId, setDepthRangeId] = useState('')
  const [color, setColor] = useState('')
  const [kind, setKind] = useState('')
  const [features, setFeatures] = useState<boolean[]>(new Array(chairSearchCondition.feature.list.length).fill(false))
  const [chairSearchParams, setChairSearchParams] = useState<ChairSearchParams | null>(null)
  const [searchResult, setSearchResult] = useState<ChairSearchResponse | null>(null)
  const [page, setPage] = useState<number>(0)

  const onSearch = () => {
    const selectedFeatures = chairSearchCondition.feature.list.filter((_, i) => features[i])
    const params: ChairSearchParams = {
      priceRangeId,
      heightRangeId,
      widthRangeId,
      depthRangeId,
      color,
      kind,
      features: selectedFeatures.length > 0 ? selectedFeatures.join(',') : '',
      page: 0,
      perPage: ESTATE_COUNTS_PER_PAGE
    }
    setChairSearchParams(params)
    setSearchResult(null)

    searchChair(params)
      .then(result => {
        setSearchResult(result)
        setPage(0)
      })
      .catch(console.error)
  }

  return (
    <>
      <Paper className={classes.page}>
        <Container maxWidth='md'>
          <Box width={1} className={classes.search}>
            <RangeForm
              name='イスの高さ'
              value={heightRangeId}
              rangeCondition={chairSearchCondition.height}
              onChange={(_, value) => { setHeightRangeId(value) }}
            />

            <RangeForm
              name='イスの横幅'
              value={widthRangeId}
              rangeCondition={chairSearchCondition.width}
              onChange={(_, value) => { setWidthRangeId(value) }}
            />

            <RangeForm
              name='イスの奥行き'
              value={depthRangeId}
              rangeCondition={chairSearchCondition.depth}
              onChange={(_, value) => { setDepthRangeId(value) }}
            />

            <RangeForm
              name='価格'
              value={priceRangeId}
              rangeCondition={chairSearchCondition.price}
              onChange={(_, value) => { setPriceRangeId(value) }}
            />

            <RadioButtonForm
              name='色'
              value={color}
              items={chairSearchCondition.color.list}
              onChange={(_, value) => { setColor(value) }}
            />

            <RadioButtonForm
              name='種類'
              value={kind}
              items={chairSearchCondition.kind.list}
              onChange={(_, value) => { setKind(value) }}
            />

            <CheckboxForm
              name='特徴'
              checkList={features}
              selectList={chairSearchCondition.feature.list}
              onChange={(_, checked, key) => {
                setFeatures(features.map((feature, i) => key === i ? checked : feature))
              }}
            />

            <Button
              onClick={onSearch}
              disabled={
                heightRangeId === '' &&
                widthRangeId === '' &&
                depthRangeId === '' &&
                priceRangeId === '' &&
                color === '' &&
                kind === '' &&
                !features.some(feature => feature)
              }
            >
              Search
            </Button>
          </Box>
        </Container>
      </Paper>

      {chairSearchParams ? (
        <Paper className={classes.page}>
          <Container maxWidth='md'>
            <Box width={1} className={classes.search} alignItems='center'>
              {searchResult ? (
                <>
                  <Pagination
                    count={Math.ceil(searchResult.count / ESTATE_COUNTS_PER_PAGE)}
                    page={page + 1}
                    onChange={(_, page) => {
                      if (!chairSearchParams) return
                      const params = { ...chairSearchParams, page: page - 1 }
                      setChairSearchParams(params)
                      setSearchResult(null)
                      const urlSearchParams = new URLSearchParams()
                      for (const [key, value] of Object.entries(params)) {
                        urlSearchParams.append(key, value.toString())
                      }
                      searchChair(params)
                        .then(result => {
                          setSearchResult(result)
                          setPage(page - 1)
                        })
                        .catch(console.error)
                    }}
                  />
                  {
                    searchResult.chairs.map((chair) => (
                      <Link key={chair.id} href={`/chair/detail?id=${chair.id}`}>
                        <Card className={classes.card}>
                          <CardActionArea className={classes.cardActionArea}>
                            <CardMedia image={chair.thumbnail} className={classes.cardMedia} />
                            <CardContent className={classes.cardContent}>
                              <h2>{chair.name}</h2>
                              <p><strong>価格:</strong> {chair.price}円</p>
                              <p><strong>詳細:</strong> {chair.description}</p>
                            </CardContent>
                          </CardActionArea>
                        </Card>
                      </Link>
                    ))
                  }
                </>
              ) : (
                <CircularProgress />
              )}
            </Box>
          </Container>
        </Paper>
      ) : null}
    </>
  )
}

const ChairSearchPage = () => {
  const [chairSearchCondition, setChairSearchCondition] = useState<ChairSearchCondition | null>(null)

  useEffect(() => {
    fetch('/api/chair/search/condition', { mode: 'cors' })
      .then(async response => await response.json())
      .then(chair => setChairSearchCondition(chair as ChairSearchCondition))
      .catch(console.error)
  }, [])

  return chairSearchCondition ? (
    <ChairSearch chairSearchCondition={chairSearchCondition} />
  ) : (
    <Loading />
  )
}

export default ChairSearchPage
