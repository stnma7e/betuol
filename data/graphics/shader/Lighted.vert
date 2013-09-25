#version 330

layout(location = 0) in vec3 position;
layout(location = 1) in vec3 normal;

smooth out vec4 interpColor;

uniform mat4 modelToCamera;
uniform mat4 projection;
uniform mat3 normalModelToCamera;

uniform vec3 lightPosition;
uniform vec4 lightIntensity;
uniform vec4 ambientIntensity;

void main()
{
	vec4 positionCameraSpace = vec4(position, 1.0) * modelToCamera;
	gl_Position = positionCameraSpace * projection;

	//mat3 normalModelToCamera = mat3(modelToCamera);
	vec3 normCamSpace = normalize(normal * normalModelToCamera);
	vec3 dirToLight = normalize(lightPosition - vec3(positionCameraSpace));

	float cosAngleIncidence = dot(normCamSpace, dirToLight);
	cosAngleIncidence = clamp(cosAngleIncidence, 0, 1);

	vec4 diffuseColor = vec4(1.0,1.0,1.0,1.0);
	interpColor = diffuseColor * (lightIntensity * cosAngleIncidence) + (diffuseColor * ambientIntensity);
	//interpColor = lightIntensity * vec4(normal, 1.0)  + (ambientIntensity);
}
