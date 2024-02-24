import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

let locales = ['en', 'ko'];

function getLocale(request: NextRequest) {
  // Check if the locale is supported
  const locale =
    request.cookies.get('locale') ||
    request.headers.get('Accept-Language') ||
    'en';
  return typeof locale === 'string' && locales.includes(locale) ? locale : 'en';
}

const PUBLIC_PATHS = ['/images', '/icons'];

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  // Check if there is any supported locale in the pathname
  const { pathname } = request.nextUrl;
  console.log(pathname);

  if (PUBLIC_PATHS.some((path) => pathname.startsWith(path)))
    return NextResponse.next();

  const pathnameHasLocale = locales.some(
    (locale) => pathname.startsWith(`/${locale}/`) || pathname === `/${locale}`
  );
  if (pathnameHasLocale) return;

  // Redirect if there is no locale
  const locale = getLocale(request);
  request.nextUrl.pathname = `/${locale}${pathname}`;
  // e.g. incoming request is /products
  // The new URL is now /en-US/products
  return Response.redirect(request.nextUrl);
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: [
    // Skip all internal paths (_next)
    '/((?!_next).*)'
    // Optional: only run on root (/) URL
    // '/'
  ]
};
