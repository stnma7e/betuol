#version 330

layout(location = 0) in vec3 position;
layout(location = 1) in vec3 normal;

smooth out vec4 interpColor;

uniform mat4 modelToCamera;
uniform mat4 projection;
uniform mat3 normalModelToCamera;

uniform vec3 dirToLight;
uniform vec4 lightIntensity;

void main()
{
	gl_Position = vec4(position, 1.0) * modelToCamera * projection;

	//mat3 normalModelToCamera = mat3(modelToCamera);
	vec3 normCamSpace = normalize(normal * normalModelToCamera);

	float cosAngleIncidence = dot(normCamSpace, vec4(dirToLight,0.0));
	cosAngleIncidence = clamp(cosAngleIncidence, 0, 1);

	vec4 color = vec4(1.0,1.0,1.0,1.0);
//	interpColor = lightIntensity * vec4(normal, 1.0);
	interpColor = lightIntensity * cosAngleIncidence;
}
