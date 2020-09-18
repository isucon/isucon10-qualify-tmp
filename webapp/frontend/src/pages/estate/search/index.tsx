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
import { CheckboxForm } from '../../../components/CheckboxForm'

import type { FC } from 'react'
import type { Estate, EstateSearchCondition, EstateSearchParams, EstateSearchResponse } from '@types'

const ESTATE_COUNTS_PER_PAGE = 20

interface EstateItemProps {
  estate: Estate
}

interface EstateSearchProps {
  estateSearchCondition: EstateSearchCondition
}

const useEstateItemStyles = makeStyles(theme =>
  createStyles({
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

const useEstateSearchStyles = makeStyles(theme =>
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
      display: 'flex',
      alignItems: 'flex-start',
      justifyContent: 'flex-start'
    },
    cardMedia: {
      width: 360,
      height: 270
    },
    cardContent: {
      marginLeft: theme.spacing(1)
    }
  })
)

const searchEstate = async (params: EstateSearchParams) => {
  const urlSearchParams = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    urlSearchParams.append(key, value.toString())
  }

  const response = await fetch(`/api/estate/search?${urlSearchParams.toString()}`, { mode: 'cors' })
  const json = await response.json()
  return json as EstateSearchResponse
}

const EstateItem: FC<EstateItemProps> = ({ estate }) => {
  const classes = useEstateItemStyles()

  return (
    <Link key={estate.id} href={`/estate/detail?id=${estate.id}`}>
      <Card className={classes.card}>
        <CardActionArea className={classes.cardActionArea}>
          <CardMedia image={estate.thumbnail} className={classes.cardMedia} />
          <CardContent className={classes.cardContent}>
            <h2>{estate.name}</h2>
            <p><strong>住所:</strong> {estate.address}</p>
            <p><strong>価格:</strong> {estate.rent}円</p>
            <p><strong>詳細:</strong> {estate.description}</p>
          </CardContent>
        </CardActionArea>
      </Card>
    </Link>
  )
}

const EstateSearch: FC<EstateSearchProps> = ({ estateSearchCondition }) => {
  const classes = useEstateSearchStyles()

  const [doorWidthRangeId, setDoorWidthRangeId] = useState('')
  const [doorHeightRangeId, setDoorHeightRangeId] = useState('')
  const [rentRangeId, setRentRangeId] = useState('')
  const [features, setFeatures] = useState<boolean[]>(new Array(estateSearchCondition.feature.list.length).fill(false))
  const [estateSearchParams, setEstateSearchParams] = useState<EstateSearchParams | null>(null)
  const [searchResult, setSearchResult] = useState<EstateSearchResponse | null>(null)
  const [page, setPage] = useState<number>(0)

  const onSearch = () => {
    const selectedFeatures = estateSearchCondition.feature.list.filter((_, i) => features[i])
    const params: EstateSearchParams = {
      doorWidthRangeId,
      doorHeightRangeId,
      rentRangeId,
      features: selectedFeatures.length > 0 ? selectedFeatures.join(',') : '',
      page: 0,
      perPage: ESTATE_COUNTS_PER_PAGE
    }
    setEstateSearchParams(params)
    setSearchResult(null)
    searchEstate(params)
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
              name='ドアの横幅'
              value={doorWidthRangeId}
              rangeCondition={estateSearchCondition.doorWidth}
              onChange={(_, value) => { setDoorWidthRangeId(value) }}
            />

            <RangeForm
              name='ドアの高さ'
              value={doorHeightRangeId}
              rangeCondition={estateSearchCondition.doorHeight}
              onChange={(_, value) => { setDoorHeightRangeId(value) }}
            />

            <RangeForm
              name='賃料'
              value={rentRangeId}
              rangeCondition={estateSearchCondition.rent}
              onChange={(_, value) => { setRentRangeId(value) }}
            />

            <CheckboxForm
              name='特徴'
              checkList={features}
              selectList={estateSearchCondition.feature.list}
              onChange={(_, checked, key) => {
                setFeatures(
                  features.map((feature, i) => key === i ? checked : feature)
                )
              }}
            />

            <Button
              onClick={onSearch}
              disabled={
                doorWidthRangeId === '' &&
                doorHeightRangeId === '' &&
                rentRangeId === '' &&
                !features.some(feature => feature)
              }
            >
              Search
            </Button>
          </Box>
        </Container>
      </Paper>

      {estateSearchParams ? (
        <Paper className={classes.page}>
          <Container maxWidth='md'>
            <Box width={1} className={classes.search} alignItems='center'>
              {searchResult ? (
                <>
                  <Pagination
                    count={Math.ceil(searchResult.count / ESTATE_COUNTS_PER_PAGE)}
                    page={page + 1}
                    onChange={(_, page) => {
                      if (!estateSearchParams) return
                      const params = { ...estateSearchParams, page: page - 1 }
                      setEstateSearchParams(params)
                      setSearchResult(null)
                      searchEstate(params)
                        .then(result => {
                          setSearchResult(result)
                          setPage(page - 1)
                        })
                        .catch(console.error)
                    }}
                  />
                  {
                    searchResult.estates.map((estate, i) => (
                      <EstateItem key={i} estate={estate} />
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

const EstateSearchPage = () => {
  const [estateSearchCondition, setEstateSearchCondition] = useState<EstateSearchCondition | null>(null)

  useEffect(() => {
    fetch('/api/estate/search/condition', { mode: 'cors' })
      .then(async response => await response.json())
      .then(estate => setEstateSearchCondition(estate as EstateSearchCondition))
      .catch(console.error)
  }, [])

  return estateSearchCondition ? (
    <EstateSearch estateSearchCondition={estateSearchCondition} />
  ) : (
    <Loading />
  )
}

export default EstateSearchPage
