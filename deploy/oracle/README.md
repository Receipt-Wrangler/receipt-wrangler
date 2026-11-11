# Oracle Cloud deployment (fork-only)

Auto-deploy for this fork onto an Oracle Cloud **VM.Standard.E2.1.Micro** (1 GB RAM, x86_64).

```
push to main ──► .github/workflows/deploy.yml ──► ghcr.io/ianrollo/receipt-wrangler:latest
                                                          ▲
                        systemd timer on the VM (5 min) ──┘ pull → restart if changed
```

GitHub holds no credentials for the server, and the server only needs to pull a public image,
so nothing secret lives in this public repo.

## 1. GitHub (once)

1. **Actions tab → enable workflows** (forks start with Actions off).
2. Disable the inherited upstream workflows, which need the upstream maintainer's secrets or burn
   minutes on e2e runs: open each of **CI - Build and Test**, **E2E**, **Mobile E2E** and
   **Release** → `···` → **Disable workflow**. Disabling them in the UI rather than deleting the
   files keeps `git merge upstream/main` conflict-free.
3. Run **Deploy - Build and Push Image** manually (**Run workflow**). The first build takes a while
   (PyTorch and Chromium layers). Later builds reuse the cache.
4. Your profile → **Packages → receipt-wrangler → Package settings**. Make sure **visibility is
   Public**, or the VM cannot pull it anonymously.

## 2. The VM (once)

SSH in as `ubuntu`.

**Swap is not optional on 1 GB.** OCR of a phone photo spikes memory well past what's free.

```bash
sudo fallocate -l 4G /swapfile && sudo chmod 600 /swapfile
sudo mkswap /swapfile && sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
echo 'vm.swappiness=10' | sudo tee /etc/sysctl.d/99-swappiness.conf && sudo sysctl --system
```

**Docker:**

```bash
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker ubuntu   # log out and back in afterwards
```

On Ubuntu 20.04 (this VM's image) the script fails at the end on `docker-model-plugin`, which
Docker doesn't publish for focal. It has already added Docker's apt repo by then, so install the
rest directly:
`sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin docker-buildx-plugin`.

**Open port 80.** Two places, both needed on Oracle:

- Console → your VCN → Security List → add an ingress rule: source `0.0.0.0/0`, TCP, port `80`.
- On the VM (Oracle's Ubuntu images ship a restrictive iptables). The rule must land *before*
  the chain's `REJECT`, whose position varies by image, so look it up rather than hard-coding it:
  ```bash
  n=$(sudo iptables -L INPUT --line-numbers -n | awk '$2=="REJECT"{print $1; exit}')
  sudo iptables -I INPUT "$n" -m state --state NEW -p tcp --dport 80 -j ACCEPT
  sudo netfilter-persistent save
  ```

**The app.** Copy this directory's files to the server, e.g. from your machine:
`scp deploy/oracle/* deploy/oracle/.env.example ubuntu@<vm-ip>:~`. Then on the VM:

```bash
sudo mkdir -p /opt/receipt-wrangler && sudo chown ubuntu: /opt/receipt-wrangler
mv ~/docker-compose.yml ~/.env.example /opt/receipt-wrangler/
cd /opt/receipt-wrangler
cp .env.example .env
sed -i "s/^SECRET_KEY=.*/SECRET_KEY=$(openssl rand -hex 32)/; s/^ENCRYPTION_KEY=.*/ENCRYPTION_KEY=$(openssl rand -hex 32)/" .env
docker compose up -d

sudo mv ~/receipt-wrangler-update.* /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now receipt-wrangler-update.timer
```

Copy the two keys out of `.env` into your password manager.

## 3. First login: do this immediately

The app is now on the public internet with the default admin **`admin` / `admin`**.

1. Log in at `http://<vm-ip>` and **change that password**.
2. **System Settings → turn off sign-ups** unless you want strangers creating accounts.

## Living within 1 GB

- **OCR engine: Tesseract.** Leave EasyOCR off. It loads PyTorch models, needs more than 1 GB
  on its own, and will be OOM-killed.
- **AI extraction:** use a hosted provider (OpenAI, Gemini, etc.). A local model will not fit.
- Expect OCR to be **slow**: the micro gets 1/8 of a core with burst. Scans process in the
  background, so the UI stays usable.
- Idle footprint is roughly 150-250 MB (API, nginx, Redis).

## Operating it

```bash
cd /opt/receipt-wrangler
docker compose logs -f app                       # app logs
systemctl list-timers receipt-wrangler-update*   # when the next update check runs
sudo systemctl start receipt-wrangler-update     # pull now instead of waiting
free -h                                          # watch swap use
```

**Backups:** everything is in `/opt/receipt-wrangler/{data,sqlite}` plus `.env`. Stop the app
(`docker compose stop app`) before copying `sqlite/` so the database isn't mid-write.

**Roll back:** every build is also tagged `sha-<commit>`. Set `image:` in `docker-compose.yml` to
that tag and run `docker compose up -d`. The timer will then keep that version until you set it back to
`latest`.

**HTTPS:** the app runs on plain HTTP. With a domain pointed at the VM, put Caddy in front for
automatic certificates, then turn on **System Settings → Secure cookies**. Turn it on only *after*
HTTPS works, or the web app locks you out.
