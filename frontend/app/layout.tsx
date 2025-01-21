import 'bootstrap/dist/css/bootstrap.min.css'
import { ReactNode } from 'react'

export const metadata = {
    title: 'Zousui Communities',
    description: 'Sample MVP with Next.js',
}

export default function RootLayout({
    children,
}: {
    children: ReactNode
}) {
    return (
        <html lang="ja">
            <body>{children}</body>
        </html>
    )
}
