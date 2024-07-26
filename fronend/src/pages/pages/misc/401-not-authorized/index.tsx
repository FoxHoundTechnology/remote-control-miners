import { ReactNode } from 'react'

import BlankLayout from 'src/@core/layouts/BlankLayout'

import Error401 from 'src/pages/401'

const Error = () => <Error401 />

Error.getLayout = (page: ReactNode) => <BlankLayout>{page}</BlankLayout>

export default Error
