// Variables for entire pipeline
local app_name = '${DRONE_REPO_NAME,,}';  // Name of app for in registry (default: repo name)
local default_docker_img = 'docker:18.09';  // Name of docker client version for building imgs
local registry = 'harbor.service.dc3.consul';  // Address of docker registry

// Functions for compiling docker commands
local docker_img(env='staging', tag='${DRONE_COMMIT}') = '%(registry)s/%(env)s/%(app_name)s/${DRONE_BRANCH,,}:%(tag)s' % { registry: registry, env: env, app_name: app_name, tag: tag };
local docker_login = 'docker login ' + registry + ' -u $${registry_user} -p $${registry_password}';
local docker_tag(tag1, tag2) = 'docker tag ' + tag1 + ' ' + tag2;
local docker_push(tag) = 'docker push ' + tag;
local docker_pull(tag) = 'docker pull ' + tag;

// Compile docker build command with many params
local docker_build_command(tag=docker_img(), target='') =

  local add_docker_param(obj) = '--%(name)s "%(value)s"' % { name: obj.name, value: obj.value };

  local default_docker_params =
    [
      { name: 'label', value: "org.opencontainers.image.created=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" },
      { name: 'label', value: 'org.opencontainers.image.version=${DRONE_COMMIT}' },
      { name: 'label', value: 'org.opencontainers.image.revision=${DRONE_COMMIT}' },
      { name: 'label', value: 'org.opencontainers.image.source=${DRONE_REPO_LINK}' },
      { name: 'label', value: 'org.opencontainers.image.title=${DRONE_REPO_NAME,,}' },
      { name: 'label', value: 'org.opencontainers.image.vendor=${DRONE_REPO_OWNER}' },
      { name: 'label', value: 'ru.qlean.drone.build.author=${DRONE_COMMIT_AUTHOR}' },
      { name: 'label', value: 'ru.qlean.drone.build.number=${DRONE_BUILD_NUMBER}' },
      { name: 'label', value: 'ru.qlean.drone.build.link=${DRONE_BUILD_LINK}' },
      { name: 'label', value: 'ru.qlean.drone.build.branch=${DRONE_BRANCH,,}' },
      { name: 'label', value: 'ru.qlean.drone.build.branch=${DRONE_BRANCH,,}' },
      { name: 'label', value: 'maintainer=${DRONE_COMMIT_AUTHOR}' },
    ];

  'DOCKER_BUILDKIT=1 docker build . ' +
  (if target != '' then std.format('--target %s ', target) else '') +
  std.format('--tag %s ', tag) +
  std.join(' ', std.map(add_docker_param, default_docker_params))
;
// Credentials for login to docker registry
local registry_credentials = {
  registry_user: { from_secret: 'registry_user' },
  registry_password: { from_secret: 'registry_password' },
};
// Default volumes for pipeline
local default_volumes = [{ name: 'docker', host: { path: '/var/run/docker.sock' } }];
// Default volume mounts for steps in pipeline
local default_mounts = [{ name: 'docker', path: '/var/run/docker.sock' }];
// Default services to run in each pipeline
local default_services = [
  {
    name: 'runner_info',
    image: default_docker_img,
    entrypoint: ['/bin/sh'],
    commands: [
      'echo ${DRONE_MACHINE}',
    ],
  },
];

local default_triggers = { event: ['push'] };

local build_and_push_pipeline(target='', tag='${DRONE_COMMIT}') = {
  name: if target != '' then 'build ' + target else 'build',
  kind: 'pipeline',
  clone: {
    depth: 1,
  },
  trigger: default_triggers,
  volumes: default_volumes,
  services: default_services,
  steps:
    [
      {
        name: 'build',
        image: default_docker_img,
        commands: [
          docker_build_command(target=target, tag=docker_img(tag=tag)),
        ],
        volumes: default_mounts,
      },
      {
        name: if target != '' then 'push built image for ' + target else 'push built image',
        image: default_docker_img,
        commands: [
          docker_login,
          docker_push(docker_img(tag=tag)),
        ],
        volumes: default_mounts,
        environment: registry_credentials,
      },
    ],
};

local release_pipeline = {
  local envs = ['production', 'staging'],
  kind: 'pipeline',
  trigger: default_triggers,
  name: 'release',
  clone: {
    disable: true,
  },
  steps: [
    {

      name: 'docker_push_' + env,
      image: default_docker_img,
      commands: [
        docker_pull(docker_img(tag='${DRONE_COMMIT}')),

        docker_tag(docker_img(tag='${DRONE_COMMIT}', env='staging'), docker_img(tag='${DRONE_COMMIT}', env=env)),
        docker_tag(docker_img(tag='${DRONE_COMMIT}', env='staging'), docker_img(tag='latest', env=env)),
        docker_tag(docker_img(tag='${DRONE_COMMIT}'), docker_img(tag='latest')),

        docker_login,
        docker_push(docker_img(tag='latest')),
        docker_push(docker_img(tag='${DRONE_COMMIT}', env=env)),
        docker_push(docker_img(tag='latest', env=env)),
      ],
      environment: registry_credentials,
      volumes: default_mounts,
      when: if env == 'production' then { branch: ['master'] } else { branch: { exclude: ['master'] } },
    }
    for env in envs
  ],
  volumes: default_volumes,
};

local slack_notification(name='', channel='drone-ci', secret='notify_slack_webhook', template='') = {
  local env = 'prod',
  name: 'slack ' + name,
  image: registry + '/library/drone-slack/master:a5c6cdfecc3a2675b1bf537232bf8d369e24ae05',
  when: { status: ['success', 'failure'] },
  settings: {
    channel: channel,
    webhook: {
      from_secret: secret,
    },
    template: if template != '' then template
    else '{{#success build.status}} *SUCCESS* {{else}} *FAIL* {{/success}} <${DRONE_BUILD_LINK}|${DRONE_REPO_NAME}#${DRONE_COMMIT:0:10}-${DRONE_BUILD_NUMBER}> {{#equal "${DRONE_BUILD_EVENT}" "pull_request"}} (PR: <${DRONE_COMMIT_LINK}|${DRONE_REPO_LINK}/pulls/${DRONE_PULL_REQUEST}>) {{/equal}} {{#equal "${DRONE_BUILD_EVENT}" "push"}} (Br: <${DRONE_REPO_LINK}/tree/${DRONE_COMMIT_BRANCH}|${DRONE_COMMIT_BRANCH}>) {{/equal}} {{#equal "${DRONE_BUILD_EVENT}" "tag"}} (Tag: <${DRONE_TAG}|${DRONE_REPO_LINK}/tree/${DRONE_TAG}>) {{/equal}} by <https://github.com/${DRONE_COMMIT_AUTHOR}|${DRONE_COMMIT_AUTHOR}>\n',
  },
};

local notification_pipeline =
  {
    kind: 'pipeline',
    trigger: default_triggers,
    name: 'notifications',
    clone: {
      disable: true,
    },
    steps: [
      slack_notification(),
    ],
  };


local deploy_pipeline(env) =
  local app = app_name;
  {
    kind: 'pipeline',
    trigger: default_triggers,
    name: 'deploy ' + env,
    clone: {
      disable: true,
    },
    steps: [
      slack_notification(
        name='deploy %(env)s start' % { env: env },
        secret='deploy_slack_webhook',
        channel='deploy',
        template='<${DRONE_BUILD_LINK}|Deploy to production> STARTED\n Project: <${DRONE_REPO_LINK}|${DRONE_REPO_NAME,,}>\n Commit: <${DRONE_COMMIT_LINK}|${DRONE_COMMIT_SHA}>\n Triggered_by: {{#equal build.event "push"}}${DRONE_COMMIT_AUTHOR}{{/equal}} {{#equal build.event "promote"}}{{build.triggeredBy}} {{/equal}}\n',
      ),
      {
        name: 'deploy',
        image: registry + '/library/infrastructure/deploy:latest',
        pull: true,
        commands: [
          'cd /infrastructure',
          '/entrypoint.sh deploy %(app)s  -e disable_parallel_deployments=true -e docker_tag=${DRONE_COMMIT} -vv' % { app: app },
        ],
        volumes: [
          {
            name: 'infrastructure',
            path: '/infrastructure',
          },
        ],
        environment: {
          CONSUL_TOKEN: { from_secret: 'deploy_consul_token' },
        },
      },
      slack_notification(
        name='deploy %(env)s finish' % { env: env },
        secret='deploy_slack_webhook',
        channel='deploy',
        template='{{#success build.status}}\n <${DRONE_BUILD_LINK}|Deploy to %(env)s> *SUCCEEDED*\n {{else}}\n <${DRONE_BUILD_LINK}|Deploy to %(env)s> *FAILED*\n {{/success}}\n Project: <${DRONE_REPO_LINK}|${DRONE_REPO_NAME,,}>\n Commit: <${DRONE_COMMIT_LINK}|${DRONE_COMMIT_SHA}>\n Triggered_by: <@{{#equal build.event "push"}}${DRONE_COMMIT_AUTHOR}{{/equal}}{{#equal build.event "promote"}}{{build.triggeredBy}}{{/equal}}>\n' % { env: env },
      ),
    ],
    volumes: [
      { name: 'infrastructure', temp: {} },
    ],
  };


local merge_dict(obj, ex_dep) =
  local merged(a) = std.mergePatch(a, ex_dep);
  if std.isArray(obj) then
    std.map(merged, obj)
  else
    merged(obj);

local set_dependencies(obj, list) =
  merge_dict(obj, { depends_on: list });

local set_trigger(obj, trigger) =
  merge_dict(obj, { trigger: trigger });

std.flattenArrays([
  [build_and_push_pipeline()],
  [set_dependencies(release_pipeline, ['build'])],
  [set_trigger(
    set_dependencies(notification_pipeline, ['release'])
    , { event: ['push'], status: ['success', 'failure'] },
  )],
  [set_trigger(
    set_dependencies(
      deploy_pipeline('production'),
      ['release'],
    )
    , { branch: 'master', event: ['promote'] },
  )],
])
