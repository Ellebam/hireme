import { getRequestConfig } from 'next-intl/server';

export default getRequestConfig(async () => {
  // For MVP, use English only
  const locale = 'en';

  return {
    locale,
    messages: (await import(`./messages/${locale}.json`)).default,
  };
});
