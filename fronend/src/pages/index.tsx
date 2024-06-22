import { useRouter } from 'next/router'
import { useEffect } from 'react'

// TODO: intial loading template

const Home = () => {
  const { push } = useRouter()

  useEffect(() => {
    push('/apps/miner/list')
  }, [])

  return <>Initial Page</>
}

export default Home
