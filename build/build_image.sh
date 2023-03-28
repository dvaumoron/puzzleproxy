#!/usr/bin/env bash

go install

buildah from --name puzzleproxy-working-container scratch
buildah copy puzzleproxy-working-container $HOME/go/bin/puzzleproxy /bin/puzzleproxy
buildah config --entrypoint '["/bin/puzzleproxy"]' puzzleproxy-working-container
buildah commit puzzleproxy-working-container puzzleproxy
buildah rm puzzleproxy-working-container

buildah push puzzleproxy docker-daemon:puzzleproxy:latest
