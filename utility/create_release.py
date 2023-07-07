import subprocess
import yaml
import sys

# take argument: version
arguments = sys.argv
version = arguments[1]
print("creating release for version: " + version)

# simple function for cmd execution


def executeCommand(cmd):
    print("executing command: " + cmd)
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode == 0:
        print("Command executed successfully!")
        print("Output:")
        print(result.stdout)
    else:
        print("Command execution failed!")
        print("Error:")
        print(result.stderr)


def loadYamlFile(path):
    with open(path, 'r') as file:
        data = yaml.safe_load(file)
    return data


def writeYamlFile(path, data):
    with open(path, 'w') as file:
        yaml.dump(data, file)


def updateHelmChartVersion():
    data = loadYamlFile("helm/Chart.yaml")
    data['version'] = version
    data['appVersion'] = version
    writeYamlFile('helm/Chart.yaml', data)


def addHelmChartToIndex():
    data = loadYamlFile("helm-release/index.yaml")
    new_entry = {
        'name': 'f2s',
        'version': version,
        'description': 'f2s deployment',
        'urls': [
            f'https://butschi84.github.io/F2S/helm-release/f2s-{version}.tgz'
        ]
    }
    data['entries']['f2s'].append(new_entry)

    writeYamlFile('helm-release/index.yaml', data)


# update version in helm cahrt
print("update version in helm chart...")
updateHelmChartVersion()

# package the helm chart
print("packaging helm chart...")
executeCommand("helm package ./helm -d ./helm-release")

# add release to index.yaml
print("add helm chart to index.yaml...")
addHelmChartToIndex()
