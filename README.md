# DevOps Lab — TaskBoard

## Цель
Пройти путь стажёр → мидл через сквозной проект: код → CI → CD → k8s → observability → инциденты.

## Ограничения
- Железо: WSL2, 2 vCPU, 7,2 ГБ для Linux, 2 ГБ swap
- Бюджет: 0 (локально), облако — только free tier позже
- Время: гибкое, спринты по 3–5 дней

## Стек (целевой)
Docker, Docker Compose, GitHub Actions, k3d, Helm, Argo CD, Terraform, Ansible,
Prometheus, Grafana, Alertmanager, Trivy, SOPS/Sealed Secrets

## Критерии успеха (общие)
- [ ] Приложение поднимается одной командой
- [ ] CI зелёный на каждый PR
- [ ] Деплой в k8s через GitOps
- [ ] Есть дашборд и рабочий алерт
- [ ] Backup/restore проверен на практике
- [ ] Написан runbook на реальный инцидент
- [ ] Каждое решение описано в ADR

## Правила работы
- Всё через Git, никаких "руками на сервере"
- Каждый спринт = README + evidence (логи, скриншоты, вывод команд)
- Ломаем осознанно, чиним по runbook
