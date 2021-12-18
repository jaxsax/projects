import '../styles/globals.css'
import React from 'react'
import type { AppProps } from 'next/app'
import { QueryClient, QueryClientProvider } from 'react-query'
import { ReactQueryDevtools } from 'react-query/devtools'
import { ToastProvider } from 'react-toast-notifications'

function MyApp({ Component, pageProps }: AppProps): JSX.Element {
    const [queryClient] = React.useState(() => new QueryClient({
        defaultOptions: {
            queries: {
            }
        },
    }));
    return (
        <QueryClientProvider client={queryClient}>
            <ToastProvider>
                <Component {...pageProps} />
            </ToastProvider>
            <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
    )
}
export default MyApp
