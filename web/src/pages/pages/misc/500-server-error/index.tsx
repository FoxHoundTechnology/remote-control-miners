import { ReactNode } from 'react'

import BlankLayout from 'src/@core/layouts/BlankLayout'

import Error500 from 'src/pages/500'

const Error = () => <Error500 />

Error.getLayout = (page: ReactNode) => <BlankLayout>{page}</BlankLayout>

export default Error
