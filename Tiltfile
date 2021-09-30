load('ext://restart_process', 'docker_build_with_restart')

local_resource(
    'gitops-bin', 
    'make bin', 
    deps=[
        './cmd', 
        './pkg',
    ]
)

docker_build_with_restart(
    'ghcr.io/weaveworks/wego-app', 
    '.',
    only=[
        './bin',
    ],
    dockerfile="dev.dockerfile",
    entrypoint='/app/build/gitops ui run',
    live_update=[
        sync('./bin', '/app/build'),
    ],
)

k8s_yaml([
    'manifests/wego-app/deployment.yaml',
    'manifests/wego-app/role-binding.yaml',
    'manifests/wego-app/role.yaml',
    'manifests/wego-app/service-account.yaml',
    'manifests/wego-app/service.yaml',
])

k8s_resource('wego-app', port_forwards='9000', resource_deps=['gitops-bin'])
