# GitHub Actions Setup Guide for Code-Factory

## 🚨 Important Note
The workflow files couldn't be pushed automatically because GitHub requires special `workflows` permission for GitHub Apps to modify workflow files. This is a security feature.

**You'll need to add these files manually through the GitHub web interface.**

---

## 📋 Step-by-Step Instructions

### Step 1: Create ci.yml

1. **Go to your repository**: https://github.com/ssdajoker/Code-Factory

2. **Navigate to Actions tab**: Click the "Actions" tab at the top of the page

3. **Set up a new workflow**:
   - If you see a "Get started with GitHub Actions" page, click **"set up a workflow yourself"**
   - OR click **"New workflow"** button, then click **"set up a workflow yourself"**

4. **Name the file**: Change the filename from `main.yml` to `ci.yml`

5. **Paste the content below** (replace everything in the editor):

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        id: setup-go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: ${{ steps.setup-go.outputs.go-cache-dir }}
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
          fail_ci_if_error: false

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run go vet
        run: go vet ./...

      - name: Run staticcheck
        uses: dominikh/staticcheck-action@v1
        with:
          version: "2023.1"

  build:
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build
        run: go build -v ./cmd/factory

      - name: Verify binary
        run: |
          ./factory --help
          ./factory version
```

6. **Commit the file**:
   - Click the green **"Commit changes..."** button (top right)
   - Add commit message: `Add CI workflow`
   - Select **"Commit directly to the main branch"** (or your preferred branch)
   - Click **"Commit changes"**

---

### Step 2: Create release.yml

1. **Go back to the Actions tab**: https://github.com/ssdajoker/Code-Factory/actions

2. **Create another new workflow**:
   - Click **"New workflow"** button
   - Click **"set up a workflow yourself"**

3. **Name the file**: Change the filename to `release.yml`

4. **Paste the content below** (replace everything in the editor):

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: v1.25.1
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

5. **Commit the file**:
   - Click the green **"Commit changes..."** button
   - Add commit message: `Add release workflow`
   - Select **"Commit directly to the main branch"** (or your preferred branch)
   - Click **"Commit changes"**

---

## ✅ Verification

After adding both files, you can verify they're set up correctly:

1. **Go to the Actions tab**: https://github.com/ssdajoker/Code-Factory/actions

2. **You should see**:
   - **CI** workflow (runs on push to main/develop branches and PRs)
   - **Release** workflow (runs when you push a tag like `v1.0.0`)

3. **The CI workflow should trigger automatically** on your next push to main or develop branch

---

## 📝 What These Workflows Do

### ci.yml - Continuous Integration
- **Triggers**: On every push to `main` or `develop` branches, and on pull requests to `main`
- **Jobs**:
  - **test**: Runs Go tests with race detection and generates coverage report
  - **lint**: Runs `go vet` and `staticcheck` for code quality
  - **build**: Builds the binary and verifies it works

### release.yml - Automated Releases
- **Triggers**: When you push a tag starting with `v` (e.g., `v1.0.0`, `v0.1.0`)
- **Jobs**:
  - **release**: Uses GoReleaser to build binaries for multiple platforms and create a GitHub release

---

## 🎯 Next Steps

1. **Add both workflow files** using the instructions above
2. **Test the CI workflow** by making a small change and pushing to your branch
3. **Create a release** by pushing a tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

---

## 🔧 Alternative: Grant Workflows Permission to GitHub App

If you want to enable automatic workflow file management in the future:

1. Go to: https://github.com/apps/abacusai/installations/select_target
2. Select your repository
3. Under "Repository permissions", find "Workflows" and set it to "Read and write"
4. Save changes

**Note**: This is optional and only needed if you want automated tools to manage your workflow files.

---

## 📞 Need Help?

If you encounter any issues:
- Check the Actions tab for error messages
- Ensure the file paths are exactly `.github/workflows/ci.yml` and `.github/workflows/release.yml`
- Verify the YAML syntax is correct (no extra spaces or tabs)

---

**Generated**: January 8, 2026
**Repository**: https://github.com/ssdajoker/Code-Factory
