export type Locale = "en" | "ru";

export type DocItem = {
  slug: string[];
  title: Record<Locale, string>;
  summary: Record<Locale, string>;
  keywords?: string[];
};

export type DocsSection = {
  id: string;
  title: Record<Locale, string>;
  items: DocItem[];
};

export type LocalizedDocItem = {
  slug: string[];
  title: string;
  summary: string;
  sectionId: string;
  sectionTitle: string;
  keywords: string[];
};

export type LocalizedDocsSection = {
  id: string;
  title: string;
  items: Array<{
    slug: string[];
    title: string;
    summary: string;
    keywords: string[];
  }>;
};

const item = (
  slug: string,
  en: string,
  ru: string,
  summaryEn: string,
  summaryRu: string,
  keywords: string[] = [],
): DocItem => ({
  slug: slug.split("/"),
  title: { en, ru },
  summary: { en: summaryEn, ru: summaryRu },
  keywords,
});

const docsSectionsData: DocsSection[] = [
  {
    id: "start",
    title: { en: "Start here", ru: "Начало работы" },
    items: [
      item(
        "introduction",
        "What is Nanami",
        "Что такое Nanami",
        "Understand what Nanami connects and protects.",
        "Узнайте, что подключает и защищает Nanami.",
        ["overview"],
      ),
      item(
        "getting-started",
        "Choose Cloud or Self-hosted",
        "Выберите Cloud или Self-hosted",
        "Choose the operating model that fits your team.",
        "Выберите подходящий вашей команде вариант запуска.",
        ["deployment", "community"],
      ),
      item(
        "quickstart",
        "First connection",
        "Первое подключение",
        "Sign in, add a device, and verify the connection.",
        "Войдите, добавьте устройство и проверьте подключение.",
        ["quickstart", "onboarding"],
      ),
    ],
  },
  {
    id: "devices",
    title: { en: "Devices", ru: "Устройства" },
    items: [
      item(
        "devices/choose-client",
        "Choose a client",
        "Выберите клиент",
        "Choose a currently supported way to connect a device.",
        "Выберите поддерживаемый способ подключения устройства.",
        ["platform", "client"],
      ),
      item(
        "ui/client-cli-quickstart",
        "Desktop and CLI",
        "Desktop и CLI",
        "Connect with an operator-provided Desktop or CLI build.",
        "Подключитесь через предоставленную оператором сборку Desktop или CLI.",
        ["desktop", "cli"],
      ),
      item(
        "devices/manage",
        "Add and manage a device",
        "Добавление и управление устройством",
        "Enroll, inspect, revoke, and recover a device.",
        "Добавьте, проверьте, отзовите или восстановите устройство.",
        ["enrollment", "offline"],
      ),
      item(
        "devices/manual-wireguard",
        "Manual WireGuard",
        "Manual WireGuard",
        "Connect an unmanaged peer without sharing its private key.",
        "Подключите неуправляемый peer без передачи приватного ключа.",
        ["wireguard", "public key"],
      ),
    ],
  },
  {
    id: "networks",
    title: { en: "Networks", ru: "Сети" },
    items: [
      item(
        "ui/managing-networks",
        "Create and manage a network",
        "Создание и управление сетью",
        "Create a private network and review its current state.",
        "Создайте приватную сеть и проверьте ее состояние.",
        ["cidr"],
      ),
      item(
        "ui/add-node-onboarding",
        "Add resources and devices",
        "Добавление ресурсов и устройств",
        "Add the right resource without exposing more than intended.",
        "Добавьте нужный ресурс без лишнего открытия доступа.",
        ["resource", "server"],
      ),
      item(
        "concepts/workspaces-groups-networks",
        "Groups and membership",
        "Группы и членство",
        "Organize people, devices, and network membership.",
        "Организуйте людей, устройства и членство в сети.",
        ["workspace", "groups"],
      ),
    ],
  },
  {
    id: "access",
    title: { en: "Access", ru: "Доступ" },
    items: [
      item(
        "access/how-access-works",
        "How access works",
        "Как работает доступ",
        "Understand subject, policy, and resource relationships.",
        "Разберитесь в связи субъекта, правила и ресурса.",
        ["default deny"],
      ),
      item(
        "access/policies",
        "Configure resource access",
        "Настройка доступа к ресурсам",
        "Use current resource editors to grant only the required access.",
        "Используйте текущие редакторы ресурсов, чтобы выдать только нужный доступ.",
        ["policy", "allow", "deny"],
      ),
      item(
        "access/roles",
        "Roles and permissions",
        "Роли и разрешения",
        "Understand workspace roles and scoped permissions.",
        "Разберитесь в ролях рабочего пространства и разрешениях.",
        ["rbac"],
      ),
      item(
        "access/explain-decision",
        "Explain an access decision",
        "Объяснение решения доступа",
        "Trace why access is allowed or denied.",
        "Проследите, почему доступ разрешен или запрещен.",
        ["explain access"],
      ),
    ],
  },
  {
    id: "services",
    title: { en: "Private Services", ru: "Приватные сервисы" },
    items: [
      item(
        "services/private-service",
        "Create a Private Service",
        "Создание Private Service",
        "Publish one exact private TCP or UDP service.",
        "Откройте один точный приватный TCP- или UDP-сервис.",
        ["connector"],
      ),
      item(
        "services/access",
        "Give users access",
        "Предоставление доступа",
        "Allow selected people or groups to reach a service.",
        "Разрешите выбранным людям или группам доступ к сервису.",
        ["service policy"],
      ),
      item(
        "services/operate",
        "Operate and troubleshoot a service",
        "Эксплуатация и диагностика сервиса",
        "Use connector-observed state to resolve availability problems.",
        "Используйте наблюдаемое состояние коннектора для диагностики.",
        ["unavailable", "connector"],
      ),
    ],
  },
  {
    id: "paths",
    title: { en: "Paths & Gateways", ru: "Пути и шлюзы" },
    items: [
      item(
        "paths/how-paths-work",
        "How paths work",
        "Как работают пути",
        "Understand direct and gateway paths.",
        "Разберитесь в прямых путях и путях через шлюз.",
        ["direct", "gateway"],
      ),
      item(
        "paths/direct-and-gateway",
        "Direct and gateway paths",
        "Прямые пути и пути через шлюз",
        "See when a connection can be direct and when it needs a gateway.",
        "Узнайте, когда соединение прямое, а когда нужен шлюз.",
        ["single hop"],
      ),
      item(
        "paths/routes",
        "Create and manage routes",
        "Создание и управление маршрутами",
        "Send selected private traffic through one eligible gateway.",
        "Направьте выбранный приватный трафик через один подходящий шлюз.",
        ["route", "cidr"],
      ),
      item(
        "paths/gateway-eligibility",
        "Gateway eligibility",
        "Доступность шлюза",
        "Understand which observed conditions make a gateway selectable.",
        "Узнайте, при каких наблюдаемых условиях шлюз доступен для выбора.",
        ["heartbeat", "online"],
      ),
      item(
        "paths/conflicts",
        "Route conflicts and unavailable paths",
        "Конфликты маршрутов и недоступные пути",
        "Resolve overlap, pending application, and unavailable-path states.",
        "Устраните пересечения, ожидание применения и недоступные пути.",
        ["conflict", "pending"],
      ),
    ],
  },
  {
    id: "dns",
    title: { en: "DNS", ru: "DNS" },
    items: [
      item(
        "dns/nameservers",
        "Configure nameservers",
        "Настройка DNS-серверов",
        "Choose the nameservers delivered to supported clients.",
        "Выберите DNS-серверы для поддерживаемых клиентов.",
        ["resolver"],
      ),
      item(
        "dns/zones-and-records",
        "Manage zones and records",
        "Управление зонами и записями",
        "Create network-scoped names without exposing raw configuration.",
        "Создавайте имена в рамках сети без раскрытия конфигурации.",
        ["zone", "record"],
      ),
      item(
        "dns/verify",
        "Verify DNS",
        "Проверка DNS",
        "Confirm that desired DNS state is applied before relying on it.",
        "Убедитесь, что заданное DNS-состояние применено.",
        ["applied", "desired"],
      ),
      item(
        "troubleshooting/dns-runtime",
        "Troubleshoot resolution",
        "Диагностика DNS",
        "Resolve names that do not update or resolve.",
        "Устраните проблемы с обновлением и разрешением имен.",
        ["stale", "lookup"],
      ),
    ],
  },
  {
    id: "control-center",
    title: { en: "Control Center", ru: "Control Center" },
    items: [
      item(
        "control-center/overview",
        "Overview",
        "Обзор",
        "Review current topology and observed status.",
        "Проверьте текущую топологию и наблюдаемое состояние.",
        ["topology"],
      ),
      item(
        "control-center/access",
        "Access",
        "Доступ",
        "Explain access from subject through policy to resource.",
        "Объясните доступ от субъекта через правило к ресурсу.",
        ["policy"],
      ),
      item(
        "control-center/routes",
        "Paths",
        "Пути",
        "Review the selected direct or gateway path.",
        "Проверьте выбранный прямой путь или путь через шлюз.",
        ["route"],
      ),
      item(
        "control-center/services",
        "Services",
        "Сервисы",
        "Review Private Services and connector-observed state.",
        "Проверьте Private Services и наблюдаемое состояние коннектора.",
        ["private service"],
      ),
      item(
        "control-center/problems",
        "Problems",
        "Проблемы",
        "Find stale, offline, and unavailable product states.",
        "Найдите устаревшие, отключенные и недоступные состояния.",
        ["reason", "error"],
      ),
    ],
  },
  {
    id: "operate",
    title: { en: "Deploy & operate", ru: "Развертывание и эксплуатация" },
    items: [
      item(
        "deployment/saas",
        "Nanami Cloud",
        "Nanami Cloud",
        "Understand the managed-service boundary.",
        "Разберитесь в границах управляемого сервиса.",
        ["saas"],
      ),
      item(
        "deployment/self-hosted",
        "Self-hosted overview",
        "Обзор Self-hosted",
        "Understand the Community topology and operator responsibilities.",
        "Разберитесь в топологии Community и ответственности оператора.",
        ["community"],
      ),
      item(
        "deployment/install-community",
        "Install Community",
        "Установка Community",
        "Check public availability before following an installation path.",
        "Проверьте публичную доступность до начала установки.",
        ["availability", "install"],
      ),
      item(
        "deployment/configuration-and-tls",
        "Configuration and TLS",
        "Конфигурация и TLS",
        "Prepare public domains, TLS, and gateway reachability.",
        "Подготовьте публичные домены, TLS и доступность шлюза.",
        ["domain", "certificate"],
      ),
      item(
        "deployment/backup-and-restore",
        "Backup and restore",
        "Резервное копирование и восстановление",
        "Protect and restore the supported Community state.",
        "Защитите и восстановите поддерживаемое состояние Community.",
        ["disaster recovery"],
      ),
      item(
        "deployment/upgrade-and-recovery",
        "Upgrade and recovery",
        "Обновление и восстановление",
        "Upgrade with a verified backup and explicit recovery boundary.",
        "Обновляйте с проверенной копией и явной границей восстановления.",
        ["rollback"],
      ),
      item(
        "deployment/community-security-checklist",
        "Security",
        "Безопасность",
        "Review the operator security responsibilities.",
        "Проверьте обязанности оператора по безопасности.",
        ["hardening"],
      ),
      item(
        "deployment/troubleshooting",
        "Troubleshooting",
        "Диагностика",
        "Start with the affected Community component and observed state.",
        "Начните с затронутого компонента Community и наблюдаемого состояния.",
        ["operator"],
      ),
    ],
  },
  {
    id: "troubleshooting",
    title: { en: "Troubleshooting", ru: "Устранение неполадок" },
    items: [
      item(
        "troubleshooting/common-issues",
        "Find a problem by symptom",
        "Поиск проблемы по симптому",
        "Choose the shortest safe diagnostic path.",
        "Выберите кратчайший безопасный путь диагностики.",
        ["help", "error"],
      ),
      item(
        "troubleshooting/auth-and-session",
        "Sign-in and MFA",
        "Вход и MFA",
        "Recover from sign-in, session, or MFA failures.",
        "Устраните ошибки входа, сессии или MFA.",
        ["login"],
      ),
      item(
        "troubleshooting/node-enrollment",
        "Device missing or offline",
        "Устройство отсутствует или не в сети",
        "Check enrollment, selected deployment, and observed status.",
        "Проверьте добавление, выбранное развертывание и состояние.",
        ["device"],
      ),
      item(
        "troubleshooting/connection-unavailable",
        "Connection unavailable",
        "Соединение недоступно",
        "Separate authorization, local runtime, and path failures.",
        "Разделите ошибки авторизации, локального runtime и пути.",
        ["connect"],
      ),
      item(
        "troubleshooting/no-eligible-gateway",
        "No eligible gateway",
        "Нет подходящего шлюза",
        "Resolve missing, stale, offline, or permission states.",
        "Устраните отсутствующее, устаревшее, offline-состояние или права.",
        ["gateway"],
      ),
      item(
        "troubleshooting/route-does-not-apply",
        "Route does not apply",
        "Маршрут не применяется",
        "Resolve desired route state that is not yet applied.",
        "Устраните заданное, но еще не примененное состояние маршрута.",
        ["route"],
      ),
      item(
        "troubleshooting/dns-does-not-resolve",
        "DNS does not resolve",
        "DNS не разрешает имя",
        "Check client settings, gateway state, and applied DNS configuration.",
        "Проверьте настройки клиента, шлюз и примененную DNS-конфигурацию.",
        ["dns"],
      ),
      item(
        "troubleshooting/private-service-unavailable",
        "Private Service unavailable",
        "Private Service недоступен",
        "Check policy, connector observation, and exact target configuration.",
        "Проверьте правило, состояние коннектора и точную конфигурацию ресурса.",
        ["service"],
      ),
      item(
        "troubleshooting/control-center-stale",
        "Stale Control Center state",
        "Устаревшее состояние Control Center",
        "Interpret stale observations without treating them as live truth.",
        "Интерпретируйте устаревшие данные, не считая их текущим состоянием.",
        ["stale"],
      ),
    ],
  },
  {
    id: "reference",
    title: { en: "Reference", ru: "Справочник" },
    items: [
      item(
        "reference/cli",
        "CLI",
        "CLI",
        "Review the supported CLI entry points.",
        "Изучите поддерживаемые команды CLI.",
        ["command"],
      ),
      item(
        "reference/authentication",
        "Authentication",
        "Аутентификация",
        "Review supported authentication methods and API credentials.",
        "Изучите поддерживаемые способы входа и учетные данные API.",
        ["mfa", "sso"],
      ),
      item(
        "reference/rbac",
        "RBAC",
        "RBAC",
        "Review customer-facing roles and permission boundaries.",
        "Изучите роли и границы разрешений для клиентов.",
        ["role"],
      ),
      item(
        "reference/endpoints",
        "API",
        "API",
        "Review the limited public API overview.",
        "Изучите ограниченный обзор публичного API.",
        ["endpoint"],
      ),
      item(
        "architecture/overview",
        "Architecture",
        "Архитектура",
        "Understand control, gateway, and encrypted data paths.",
        "Разберитесь в управлении, шлюзах и зашифрованных путях данных.",
        ["wireguard"],
      ),
      item(
        "reference/availability",
        "Platform and feature availability",
        "Доступность платформ и функций",
        "Check what is supported, externally gated, or unavailable.",
        "Проверьте, что поддерживается, требует внешней готовности или недоступно.",
        ["platform"],
      ),
      item(
        "reference/errors",
        "Error and reason codes",
        "Коды ошибок и причин",
        "Use public-safe codes to choose a recovery action.",
        "Используйте безопасные коды для выбора действия.",
        ["reason code"],
      ),
    ],
  },
];

function dedupeDocItems<T extends { slug: string[] }>(items: T[]) {
  const seen = new Set<string>();
  return items.filter((entry) => {
    const key = entry.slug.join("/");
    if (seen.has(key)) return false;
    seen.add(key);
    return true;
  });
}

export const allDocs = dedupeDocItems(
  docsSectionsData.flatMap((section) => section.items),
);

export const getDocPath = (slug: string[]) => `/${slug.join("/")}`;

export function getDocsSections(locale: Locale): LocalizedDocsSection[] {
  return docsSectionsData.map((section) => ({
    id: section.id,
    title: section.title[locale] ?? section.title.en,
    items: section.items.map((entry) => ({
      slug: entry.slug,
      title: entry.title[locale] ?? entry.title.en,
      summary: entry.summary[locale] ?? entry.summary.en,
      keywords: entry.keywords ?? [],
    })),
  }));
}

function matchesQuery(entry: LocalizedDocItem, query: string) {
  return [entry.title, entry.summary, entry.sectionTitle, ...entry.keywords]
    .join(" ")
    .toLowerCase()
    .includes(query);
}

export function filterDocsSections(locale: Locale, query: string) {
  const trimmed = query.trim().toLowerCase();
  const sections = getDocsSections(locale);
  if (!trimmed) return sections;

  return sections
    .map((section) => ({
      ...section,
      items: section.items.filter((entry) =>
        matchesQuery(
          {
            ...entry,
            sectionId: section.id,
            sectionTitle: section.title,
          },
          trimmed,
        ),
      ),
    }))
    .filter((section) => section.items.length > 0);
}

export function getLocalizedDocs(locale: Locale): LocalizedDocItem[] {
  return dedupeDocItems(
    docsSectionsData.flatMap((section) =>
      section.items.map((entry) => ({
        slug: entry.slug,
        title: entry.title[locale] ?? entry.title.en,
        summary: entry.summary[locale] ?? entry.summary.en,
        sectionId: section.id,
        sectionTitle: section.title[locale] ?? section.title.en,
        keywords: entry.keywords ?? [],
      })),
    ),
  );
}

export function getDocNavigation(slug: string[], locale: Locale) {
  const docs = getLocalizedDocs(locale);
  const currentIndex = docs.findIndex(
    (entry) => entry.slug.join("/") === slug.join("/"),
  );
  if (currentIndex === -1) return { current: null, previous: null, next: null };

  return {
    current: docs[currentIndex] ?? null,
    previous: currentIndex > 0 ? docs[currentIndex - 1] : null,
    next: currentIndex < docs.length - 1 ? docs[currentIndex + 1] : null,
  };
}

export function getDocTitle(slug: string[], locale: Locale = "en") {
  const found = allDocs.find(
    (entry) => entry.slug.join("/") === slug.join("/"),
  );
  if (!found) return locale === "ru" ? "Документация" : "Documentation";
  return found.title[locale] ?? found.title.en;
}
