# EduHub Backend

EduHub Backend — backend-часть CRM-системы для частных школ, образовательных центров и компаний с несколькими филиалами.

Проект строится как multi-tenant SaaS-система. Это значит, что одна платформа может обслуживать несколько компаний, но данные каждой компании изолированы через organization_id, а данные филиалов — через branch_id.

## Основная идея

Система нужна для управления:

- компаниями и школами;
- филиалами;
- пользователями;
- ролями и правами;
- преподавателями;
- учениками;
- родителями;
- предметами;
- группами;
- расписанием;
- занятиями;
- посещаемостью;
- домашними заданиями;
- оплатами;
- зарплатами преподавателей;
- файлами;
- уведомлениями;
- аудитом действий.

## Стек технологий

- Go
- PostgreSQL
- Redis
- Docker Compose
- JWT
- REST API
- golang-migrate
- chi router
- pgx

## Multi-tenant архитектура

В проекте используется одна база данных PostgreSQL.

Почти все основные таблицы имеют поле organization_id.  
Таблицы, связанные с филиалами, также имеют branch_id.

Это нужно, чтобы одна компания не могла видеть данные другой компании.

## Структура проекта

cmd/api — точка входа приложения  
configs — конфигурационные файлы  
docs — документация API  
migrations — SQL-миграции  
scripts — вспомогательные скрипты  
tests — тесты  

internal/app — приложение и роутер  
internal/config — загрузка .env  
internal/platform — база, Redis, JWT, storage  
internal/middleware — middleware  
internal/modules — бизнес-модули  
internal/shared — общие helper-пакеты  

## Основные модули

auth — авторизация  
organization — компании  
branch — филиалы  
user — пользователи  
role — роли  
permission — права доступа  
subject — предметы  
teacher — преподаватели  
student — ученики  
parent — родители  
group — группы  
schedule — расписание  
lesson — занятия  
attendance — посещаемость  
homework — домашние задания  
payment — оплаты  
payroll — зарплата преподавателей  
notification — уведомления  
file — файлы  
audit — аудит действий  

## Зарплата преподавателей

В проекте предусмотрена гибкая система расчёта зарплаты.

Админ компании сможет задавать:

- оплату за урок;
- длительность урока;
- формулу расчёта зарплаты;
- правила по филиалу;
- правила по предмету;
- отдельные правила для преподавателя.

Пример формулы:

lessons_count * lesson_rate + bonus - penalty

Важно: формулы должны выполняться безопасно на backend-стороне. Нельзя выполнять формулу как обычный код или SQL.

## Локальный запуск

1. Запустить PostgreSQL и Redis:

docker compose up -d

2. Применить миграции:

migrate -path migrations -database "postgres://eduhub:eduhub_password@localhost:5432/eduhub?sslmode=disable" up

3. Запустить backend:

go run ./cmd/api

4. Проверить сервер:

http://localhost:8080/health

5. Проверить базу:

http://localhost:8080/health/db

## Полезные команды

go test ./...

go mod tidy

gofmt -w .

docker compose up -d

docker compose down

## Текущее состояние

На данный момент готово:

- базовая структура Go-проекта;
- Docker Compose с PostgreSQL и Redis;
- загрузка конфигурации из .env;
- endpoint /health;
- endpoint /health/db;
- подключение к PostgreSQL;
- SQL-миграции;
- multi-tenant структура базы данных;
- таблицы пользователей, ролей и прав;
- таблицы учебного процесса;
- таблицы оплат и зарплат;
- payroll rules;
- таблицы файлов, уведомлений и аудита;
- базовые permissions;
- индексы для основных таблиц.

## Что сделать дальше

1. Регистрация компании и первого администратора.
2. Login, refresh token, logout и /me.
3. JWT middleware.
4. Tenant middleware.
5. RBAC / permission middleware.
6. CRUD для филиалов.
7. CRUD для предметов.
8. CRUD для преподавателей.
9. CRUD для учеников и родителей.
10. Группы, расписание и занятия.
11. Посещаемость.
12. Оплаты и зарплаты.
13. Безопасный калькулятор payroll-формул.
14. Swagger-документация.
15. Тесты.

## Статус

Проект находится на этапе начальной backend-разработки.  
База данных и базовая инфраструктура уже подготовлены.

## License

Private project.
